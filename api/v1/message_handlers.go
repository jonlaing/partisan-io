package v1

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/db"
	m "partisan/models"
	"strconv"
	"time"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"partisan/Godeps/_workspace/src/github.com/gorilla/websocket"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

type ThreadResp struct {
	ThreadUser m.MessageThreadUser `json:"thread_user"`
	HasUnread  bool                `json:"has_read"`
}

func MessageThreadIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	threads, err := dao.GetMessageThreadUsers(user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var tresps []ThreadResp
	for _, t := range threads {
		// Possible N+1 problem, but there shouldn't be so manythreads that it's an issue... maybe...
		// Not handling error, because the worst that happens is that a thread looks unread when it isn't.
		// We'll live...
		hasUnread, _ := dao.MessageThreadHasUnread(t.ThreadID, db)
		tresps = append(tresps, ThreadResp{t, hasUnread})
	}

	c.JSON(http.StatusOK, gin.H{"threads": tresps})
}

func MessageThreadCreate(c *gin.Context) {
	db := db.GetDB(c)

	currentUser, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	uID := c.Request.FormValue("user_id")
	if len(uID) == 0 {
		c.AbortWithError(http.StatusNotAcceptable, errors.New("No userID specified"))
		return
	}

	userID, err := strconv.Atoi(uID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// can only start a thread with friends
	_, err = dao.GetFriendship(currentUser, uint64(userID), db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	thread, err := dao.GetMessageThreadByUsers(currentUser.ID, uint64(userID), db)
	if err == nil {
		// Apparently, this thread already exists
		c.JSON(http.StatusOK, gin.H{"thread": thread})
		return
	}

	if _, ok := err.(*dao.MessageThreadUnreciprocated); ok {
		// somehow we didn't create all the MessageThreadUsers
		// last time, so we have to correct that
		var mtu m.MessageThreadUser
		if err := db.Where("thread_id = ? AND user_id IN (?)", []uint64{currentUser.ID, uint64(userID)}).First(&mtu).Error; err == nil {
			if mtu.UserID == currentUser.ID {
				db.Create(m.MessageThreadUser{ThreadID: thread.ID, UserID: uint64(userID)})
			} else {
				db.Create(m.MessageThreadUser{ThreadID: thread.ID, UserID: currentUser.ID})
			}
		}
	} else {
		// unknown error
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	thread = m.MessageThread{}
	if err := db.Create(&thread).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	mtu1 := m.MessageThreadUser{UserID: currentUser.ID, ThreadID: thread.ID}
	mtu2 := m.MessageThreadUser{UserID: uint64(userID), ThreadID: thread.ID}

	db.Create(&mtu1)
	db.Create(&mtu2)

	c.JSON(http.StatusOK, gin.H{"thread": thread})
}

func MessageIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	tID := c.Param("thread_id")
	if len(tID) == 0 {
		c.AbortWithError(http.StatusNotAcceptable, errors.New("No thread specified"))
		return
	}

	threadID, err := strconv.Atoi(tID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if hasUser, err := dao.MessageThreadHasUser(user.ID, uint64(threadID), db); err != nil || !hasUser {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// will also attach m.User to each m.Message
	msgs, err := dao.GetMessages(uint64(threadID), db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	dao.MarkAllMessagesRead(user.ID, uint64(threadID), db)

	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

func MessageCreate(c *gin.Context) {
	db := db.GetDB(c)
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	tID := c.Param("thread_id")
	if len(tID) == 0 {
		c.AbortWithError(http.StatusNotAcceptable, errors.New("No thread specified"))
		return
	}

	threadID, err := strconv.Atoi(tID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	thread, err := dao.GetMessageThread(uint64(threadID), db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if hasUser, err := dao.MessageThreadHasUser(user.ID, thread.ID, db); err != nil || !hasUser {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	msg := m.Message{
		UserID:    user.ID,
		ThreadID:  thread.ID,
		Body:      c.Request.FormValue("body"),
		Read:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&msg).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	thread.UpdatedAt = time.Now()
	db.Save(&thread) // Not handling error here, because it's really not that big a deal

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func MessageCount(c *gin.Context) {
	db := db.GetDB(c)
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	msg := make(chan bool)
	quit := make(chan bool)

	go messageCountReadLoop(conn, msg, quit)
	go messageCountWriteLoop(user.ID, db, conn, msg, quit)
}

func MessageSocket(c *gin.Context) {
	db := db.GetDB(c)
	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	tID := c.Param("thread_id")
	if len(tID) == 0 {
		c.AbortWithError(http.StatusNotAcceptable, errors.New("No thread specified"))
		return
	}

	threadID, err := strconv.Atoi(tID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if hasUser, err := dao.MessageThreadHasUser(user.ID, uint64(threadID), db); err != nil || !hasUser {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	msg := make(chan string)
	quit := make(chan bool)

	go messageReadLoop(conn, msg, quit)
	go messageWriteLoop(user.ID, uint64(threadID), db, conn, msg, quit)
}

func messageCountReadLoop(c *websocket.Conn, send chan bool, quit chan bool) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			quit <- true
			return
		}

		send <- true
	}
}

func messageCountWriteLoop(userID uint64, db *gorm.DB, c *websocket.Conn, received chan bool, quit chan bool) {
	for {
		select {
		case <-received:
			count, err := dao.MessageUnreadCount(userID, db)
			if err != nil {
				return
			}

			if err := c.WriteJSON(gin.H{"count": count}); err != nil {
				return
			}
		case <-quit:
			return // kill the loop
		}
	}
}

func messageReadLoop(c *websocket.Conn, send chan string, quit chan bool) {
	for {
		sMsgType, r, err := c.NextReader()
		if err != nil {
			c.Close()
			quit <- true
			return
		}

		if sMsgType == websocket.BinaryMessage {
			fmt.Println("MessageSocket message shouldn't be a binary")
		}

		stamp, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println("couldn't read timestamp", r)
		}

		send <- string(stamp)
	}
}

func messageWriteLoop(userID, threadID uint64, db *gorm.DB, c *websocket.Conn, received chan string, quit chan bool) {
	for {
		select {
		case stamp := <-received:
			sec, err := strconv.Atoi(stamp)
			if err != nil {
				fmt.Println("bad timestamp:", string(stamp))
			}

			after := time.Unix(int64(sec), int64(0))

			msgs, err := dao.GetMessagesAfter(threadID, after, db)
			if err != nil {
				c.WriteJSON(gin.H{"messages": []m.Message{}})
			}

			if err := c.WriteJSON(gin.H{"messages": msgs}); err != nil {
				return
			}

			dao.MarkAllMessagesRead(userID, uint64(threadID), db)
		case <-quit:
			return // kill the loop
		}
	}
}
