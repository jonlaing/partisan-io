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
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// NotificationsIndex shows most recent unread notifications, or the most recent 10, whichever is bigger
func NotificationsIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, err := auth.CurrentUser(c, &db)
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
					resp = append(resp, NotifResp{Notification: n, User: u})
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
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	NotificationWebsocket(c, &db)

}

// NotificationsRead sets a notification as "seen"
func NotificationsRead(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	user, err := auth.CurrentUser(c, &db)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	db.Where("target_user_id = ?", user.ID).Update("seen", true)

	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}

// NotificationWebsocket gets called in NotificationCount and returns a socket to allow us to continually poll notifications
func NotificationWebsocket(c *gin.Context, db *gorm.DB) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		user, err := auth.CurrentUser(c, db)
		if err != nil {
			conn.WriteJSON(gin.H{"error": "unauthorized"})
			return
		}

		var count int
		db.Model(m.Notification{}).Where("target_user_id = ? AND seen = ?", user.ID, false).Count(&count)

		conn.WriteJSON(gin.H{"count": count})
	}
}
