package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"github.com/gin-gonic/gin"
	"partisan/models.v2/tickets"
)

func SocketTicketCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if t, err := tickets.GetByUserID(user.ID, db); err == nil {
		if t.IsValid() {
			c.JSON(http.StatusOK, gin.H{"ticket": t})
			return
		}

		db.Delete(&t)
	}

	t := tickets.NewSocketTicket(user.ID)
	if err := db.Save(&t).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": t})
}
