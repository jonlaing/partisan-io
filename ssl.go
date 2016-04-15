package main

import (
	"net/http"
	"os"
	"partisan/logger"

	"github.com/gin-gonic/gin"
)

func ForceSSL(c *gin.Context) {
	logger.Trace.Println(c.Request)

	if os.Getenv("GIN_MODE") != "release" {
		logger.Warning.Println("Not in \"release\" mode, ignoring SSL")
		c.Next()
		return
	}

	if c.Request.Header.Get("X-Forwarded-Proto") != "https" {
		logger.Warning.Println("Trying to access site without SSL, redirecting")
		c.Redirect(http.StatusMovedPermanently, "https://www.partisan.io"+c.Request.URL.String())
		return
	}

	c.Next()
}
