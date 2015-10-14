package main

import (
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"

	"partisan/Godeps/_workspace/src/github.com/DeanThompson/ginpprof"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.Wrap(router)

	router.Run(":8080")
}
