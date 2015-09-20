package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

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

	c.JSON(http.StatusOK, notifs)
}

// NotificationsCount returns the number of unread notifications
func NotificationsCount(c *gin.Context) {
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

	var count int
	db.Model(m.Notification{}).Where("target_user_id = ? AND seen = ?", user.ID, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// NotificationsRead sets a notification as "seen"
func NotificationsRead(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	id := c.Param("record_id")

	db.Where("id = ?", id).Update("seen", true)

	c.JSON(http.StatusOK, gin.H{"message": "marked read"})
}
