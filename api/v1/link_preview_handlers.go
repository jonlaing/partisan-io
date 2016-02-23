package v1

// func LinkPreviewShow(c *gin.Context) {
// 	user, _ := auth.CurrentUser(c)

// 	link := c.Param("link")

// 	response, err := http.Get(link)
// 	if err != nil {
// 		c.AbortWithError(http.StatusNotFound, err)
// 		return
// 	}
// 	defer response.Body.Close()

// 	meta := htmlmeta.Extract(response.Body)

// 	// Try to find malicious stuff in the image
// 	// if(len(htmlmeta.OGImage) > 0) {
// 	// 	// TODO: FIGURE THIS SHIT OUT!
// 	// }
// 	c.JSON(http.StatusOK, meta)
// }
