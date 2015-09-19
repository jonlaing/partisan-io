package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"partisan/db"
	m "partisan/models"
)

// ImageAttachmentIndex gets all the attachments associated with a record
func ImageAttachmentIndex(c *gin.Context) {
	db, err := db.InitDB()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	var attachments []m.ImageAttachment

	rID, rType, err := getRecord(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if err := db.Where("record_type = ? AND record_id = ?", rType, rID).Find(&attachments).Error; err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, attachments)
	return
}
