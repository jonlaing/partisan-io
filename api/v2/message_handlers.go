package v2

import (
	"errors"
	"log"
	"net/http"
	"partisan/auth"
	"partisan/db"
	"partisan/logger"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"partisan/models.v2/messages"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var msgCWsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ThreadIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	threads, err := messages.ListThreads(user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	log.Println(len(threads))

	c.JSON(http.StatusOK, gin.H{"threads": threads})
}

func ThreadCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var binding messages.ThreadCreatorBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	thread, errs := messages.NewThread(user.ID, binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	existing, err := messages.GetByUsers(db, thread.Users.GetUserIDs()...)
	if err == nil {
		// apparently this thread already exists
		c.JSON(http.StatusOK, gin.H{"thread": existing})
		return
	}

	if err := db.Create(&thread).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	for i := range thread.Users {
		if err := db.Create(&thread.Users[i]).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"thread": thread})

}

func MessageIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	threadID := c.Param("thread_id")

	if hasUser, err := messages.HasUser(user.ID, threadID, db); err != nil || !hasUser {
		c.AbortWithError(http.StatusUnauthorized, messages.ErrThreadUser)
		return
	}

	// will also attach m.User to each m.Message
	msgs, err := messages.GetMessages(threadID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := messages.MarkAllMessagesRead(user.ID, threadID, db); err != nil {
		logger.Error.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

func NewMessages(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	threadID := c.Param("thread_id")

	if hasUser, err := messages.HasUser(user.ID, threadID, db); err != nil || !hasUser {
		c.AbortWithError(http.StatusUnauthorized, messages.ErrThreadUser)
		return
	}

	stamp, ok := c.GetQuery("timestamp")
	if !ok {
		c.AbortWithError(http.StatusNotAcceptable, errors.New("Bad timestamp"))
		return
	}

	sec, err := strconv.Atoi(stamp)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, errors.New("Bad timestamp"))
		return
	}

	after := time.Unix(int64(sec), int64(0))

	msgs, err := messages.GetMessagesAfter(threadID, after, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err := messages.MarkAllMessagesRead(user.ID, threadID, db); err != nil {
		logger.Error.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

func MessageCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	threadID := c.Param("thread_id")

	thread, err := messages.GetThread(threadID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if hasUser, err := messages.HasUser(user.ID, thread.ID, db); err != nil || !hasUser {
		c.AbortWithError(http.StatusUnauthorized, messages.ErrThreadUser)
		return
	}

	var binding messages.MessageCreatorBinding
	if err := c.BindJSON(&binding); err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	binding.UserID = user.ID
	binding.ThreadID = thread.ID

	msg, errs := messages.NewMessage(binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Create(&msg).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": msg})

	// Touch Updated at on thread and thread users
	// Not handling error here, because it's really not that big a deal
	thread.UpdatedAt = time.Now()
	db.Save(&thread)
	db.Model(messages.ThreadUser{}).Where("thread_id = ?", thread.ID).UpdateColumn("updated_at", time.Now())
}

func MessageUnread(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	conn, err := msgCWsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	msg := make(chan bool)
	quit := make(chan bool)

	go msgReadLoop(conn, msg, quit)
	go msgWriteLoop(user.ID, db, conn, msg, quit)
}

func msgReadLoop(c *websocket.Conn, send chan bool, quit chan bool) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			quit <- true
			return
		}

		send <- true
	}
}

func msgWriteLoop(userID string, db *gorm.DB, c *websocket.Conn, received chan bool, quit chan bool) {
	for {
		select {
		case <-received:
			count, err := messages.UnreadCount(userID, db)
			if err != nil {
				log.Println(err)
			}

			if err := c.WriteJSON(gin.H{"unread": count > 0}); err != nil {
				log.Println(err)
			}
		case <-quit:
			return // kill the loop
		}
	}
}
