package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/messages"

	"github.com/gin-gonic/gin"
)

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

	thread, errs := messages.NewThread(binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	existing, err := messages.GetByUsers(db, thread.Users.GetUserIDs()...)
	if err == nil {
		// apparently this thread already exists
		c.JSON(http.StatusOK, gin.H{"thread": thread})
		return
	}

	if err == messages.ErrThreadUnreciprocated {
		for _, u := range existing.Users {
			db.Delete(&u) // remove all unreciprocated users
		}
	}

	if err := db.Create(&thread).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	for i := range thread.Users {
		if err := db.Create(&thread.Users[i]).Error; err != nil {
			c.AborthWithError(http.StatusNotAcceptable, err)
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"thread": thread})

}
