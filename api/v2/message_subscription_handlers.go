package v2

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/messages"
	"partisan/models.v2/users"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

var messageWsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// MessageThreadSubscription manages a subscription
type MessageThreadSubscription struct {
	conn   *websocket.Conn
	thread messages.Thread
	user   users.User
	msg    chan string
	quit   chan bool
}

// MessageThreadSubscribe manages a websocket subscription to check for incoming messages
func MessageThreadSubscribe(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	threadID := c.Param("thread_id")

	mt, err := messages.GetThread(threadID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if ok, err := messages.HasUser(user.ID, mt.ID, db); !ok || err != nil {
		c.AbortWithError(http.StatusUnauthorized, messages.ErrThreadUser)
		return
	}

	conn, err := messageWsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}

	subscription := MessageThreadSubscription{
		conn:   conn,
		thread: mt,
		user:   user,
		msg:    make(chan string),
		quit:   make(chan bool),
	}

	go subscription.MessageReadLoop()
	go subscription.MessageWriteLoop(db)
}

func (sub MessageThreadSubscription) MessageReadLoop() {
	for {
		sMsgType, r, err := sub.conn.NextReader()
		if err != nil {
			sub.conn.Close()
			sub.quit <- true
			return
		}

		if sMsgType == websocket.BinaryMessage {
			log.Println("MessageSocket message shouldn't be a binary")
		}

		stamp, err := ioutil.ReadAll(r)
		if err != nil {
			log.Println("couldn't read timestamp", r)
		}

		sub.msg <- string(stamp)
	}
}

func (sub MessageThreadSubscription) MessageWriteLoop(db *gorm.DB) {
	for {
		select {
		case stamp := <-sub.msg:
			sec, err := strconv.Atoi(stamp)
			if err != nil {
				log.Println("bad timestamp:", string(stamp))
			}

			after := time.Unix(int64(sec), int64(0))

			msgs, err := messages.GetMessagesAfter(sub.thread.ID, after, db)
			if err != nil {
				sub.conn.WriteJSON(gin.H{"messages": []messages.Message{}})
			}

			if err := sub.conn.WriteJSON(gin.H{"messages": msgs}); err != nil {
				return
			}

			err = messages.MarkAllMessagesRead(sub.user.ID, sub.thread.ID, db)
			if err != nil {
				log.Println("Error marking messages read:", sub.user.ID, sub.thread.ID, err)
			}
		case <-sub.quit:
			return // kill the loop
		}
	}
}
