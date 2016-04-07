package v1

import (
	"net/http"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"

	"github.com/gin-gonic/gin"
)

// SocketTicketCreate creates a ticket for a mobile user to make websocket requests
func SocketTicketCreate(c *gin.Context) {
	db := db.GetDB(c)
	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	var ticket m.SocketTicket
	if err := db.Where("user_id = ?", user.ID).Find(&ticket).Error; err == nil {
		// if the ticket already exists, and it's not expired, go forth
		if ticket.IsValid() {
			c.JSON(http.StatusOK, ticket)
			return
		}

		db.Delete(&ticket) // if it's not valid, get rid of it and create a new one
	}

	ticket, err = m.NewSocketTicket(user.ID)
	if err != nil {
		handleError(err, c)
	}

	if err := db.Create(&ticket).Error; err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, ticket)
}
