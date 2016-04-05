package v1

import (
	"net/http"
	"partisan/db"
	m "partisan/models"

	"github.com/gin-gonic/gin"
)

// ImageAttachmentIndex gets all the attachments associated with a record
func ImageAttachmentIndex(c *gin.Context) {
	db := db.GetDB(c)

	var attachments []m.ImageAttachment

	rID, rType, err := getRecord(c)
	if err != nil {
		handleError(err, c)
		return
	}

	if err := db.Where("record_type = ? AND record_id = ?", rType, rID).Find(&attachments).Error; err != nil {
		handleError(&ErrDBNotFound{err}, c)
		return
	}

	c.JSON(http.StatusOK, attachments)
	return
}
