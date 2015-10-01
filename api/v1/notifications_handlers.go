package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

// NotifResp is the response for the Notification Index
type NotifResp struct {
	Notification m.Notification `json:"notification"`
	User         m.User         `json:"user"`
	Record       interface{}    `json:"record"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NotificationsIndex shows most recent unread notifications, or the most recent 10, whichever is bigger
func NotificationsIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var notifs, readNotifs []m.Notification
	var unreadCount int
	if err := db.Where("target_user_id = ? AND seen = ?", user.ID, false).Order("created_at desc").Find(&notifs).Count(&unreadCount).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if unreadCount > 0 && unreadCount < 10 {
		db.Where("target_user_id = ? AND seen = ?", user.ID, true).Order("created_at desc").Limit(10 - unreadCount).Find(&readNotifs)
		notifs = append(notifs, readNotifs...)
	}

	if len(notifs) > 0 {
		var userIDs []uint64
		for _, n := range notifs {
			userIDs = append(userIDs, n.UserID)
		}

		var users []m.User
		db.Where("id IN (?)", userIDs).Find(&users)

		// only looping through <= 10 of each, so this should be pretty performant
		var resp []NotifResp
		for _, n := range notifs {
			for _, u := range users {
				if u.ID == n.UserID {
					var r interface{}
					r, err := n.GetRecord(db)
					if err != nil {
						fmt.Println(err)
					}
					resp = append(resp, NotifResp{Notification: n, User: u, Record: r})
				}
			}
		}

		db.Model(m.Notification{}).Where("target_user_id = ?", user.ID).Update("seen", true)

		c.JSON(http.StatusOK, resp)
		return
	}

	c.JSON(http.StatusOK, []NotifResp{})
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
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	msg := make(chan bool)
	quit := make(chan bool)

	go readLoop(conn, msg, quit)
	go writeLoop(user.ID, db, conn, msg, quit)

}

// NotificationsRead sets a notification as "seen"
func NotificationsRead(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	db.Where("target_user_id = ?", user.ID).Update("seen", true)

	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}

func readLoop(c *websocket.Conn, send chan bool, quit chan bool) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			quit <- true
			return
		} else {
			send <- true
		}
	}
}

func writeLoop(userID uint64, db *gorm.DB, c *websocket.Conn, received chan bool, quit chan bool) {
	var count int

	for {
		select {
		case <-received:
			if err := db.Table("notifications").Where("target_user_id = ? AND seen = ?", userID, false).Count(&count).Error; err != nil {
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
