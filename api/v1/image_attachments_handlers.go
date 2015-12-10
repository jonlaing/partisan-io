package v1

import (
	"net/http"
	"partisan/db"
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// ImageAttachmentIndex gets all the attachments associated with a record
func ImageAttachmentIndex(c *gin.Context) {
	db := db.GetDB(c)

	var attachments []m.ImageAttachment

	rID, rType, err := getRecord(c)
	if err != nil {
		return handleError(err, c)
	}

	if err := db.Where("record_type = ? AND record_id = ?", rType, rID).Find(&attachments).Error; err != nil {
		return handleError(&ErrDBNotFound{err}, c)
	}

	c.JSON(http.StatusOK, attachments)
	return
}
