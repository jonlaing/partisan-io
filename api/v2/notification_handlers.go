package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"github.com/gorilla/websocket"

	"partisan/models.v2/notifications"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NotificationIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	notifs, err := notifications.ListByUserID(user.ID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	notifs.MarkRead(db)

	c.JSON(http.StatusOK, gin.H{"notifications": notifs})
}

// NotificationsCount returns the number of unread notifications
func NotificationsCount(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	msg := make(chan bool)
	quit := make(chan bool)

	go notifReadLoop(conn, msg, quit)
	go notifWriteLoop(user.ID, db, conn, msg, quit)

}

func notifReadLoop(c *websocket.Conn, send chan bool, quit chan bool) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			quit <- true
			return
		}

		send <- true
	}
}

func notifWriteLoop(userID string, db *gorm.DB, c *websocket.Conn, received chan bool, quit chan bool) {
	for {
		select {
		case <-received:
			count, err := notifications.CountByUserID(userID, db)
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
