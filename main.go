package main

import (
	"log"
	"net/http"
	"os"
	"partisan/db"
	"partisan/logger"
	"time"

	"partisan/models.v2/apptokens"
	"partisan/models.v2/attachments"
	"partisan/models.v2/events"
	"partisan/models.v2/flags"
	"partisan/models.v2/friendships"
	"partisan/models.v2/hashtags"
	"partisan/models.v2/messages"
	"partisan/models.v2/notifications"
	"partisan/models.v2/posts"
	"partisan/models.v2/tickets"
	"partisan/models.v2/users"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func init() {
	// apiV1.ConfigureEmailer(emailConfig)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	r := gin.Default()
	r.Use(db.DB())
	r.Use(ForceSSL)
	// recover from panics with a 500
	r.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				logger.Error.Println("Recovered Panic (main.go):", r)
			}
		}()

		c.Next()
	})
	// r.Us (gin.BasicAuth(gin.Accounts{
	// 	"partisan-basic": "antistate123",
	// }))

	// initRoutesV1(r)
	initRoutesV2(r)

	r.Use(static.Serve("/localfiles", static.LocalFile("localfiles", false)))
	r.Use(static.Serve("/", static.LocalFile("front_dist", false)))

	// homepage
	r.GET("/", func(c *gin.Context) {
		sess := sessions.Default(c)

		if sess.Get("user_id") != nil {
			c.Redirect(http.StatusFound, "/feed")
			return
		}

		c.File("front_dist/index.html")
	})

	// DON'T DO THIS IN PROD!!!
	db.Database.AutoMigrate(
		&apptokens.AppToken{},
		&posts.Post{},
		&users.User{},
		&friendships.Friendship{},
		&attachments.Attachment{},
		&notifications.Notification{},
		&hashtags.Hashtag{},
		&hashtags.Taxonomy{},
		&flags.Flag{},
		// &m.UserTag{},
		&messages.Message{},
		&messages.Thread{},
		&messages.ThreadUser{},
		&tickets.SocketTicket{},
		&events.Event{},
		&events.EventSubscription{},
	)

	ginpprof.Wrapper(r)

	// r.Run(":" + os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   3 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

// func deprecated(newPath string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if newPath != "" {
// 			c.AbortWithError(http.StatusBadRequest, fmt.Errorln("This endpoint is deprecated. Please see:", newPath))
// 			return
// 		}

// 		c.AbortWithError(http.StatusBadRequest, fmt.Errorln("This endpoint is deprecated."))
// 	}
// }
