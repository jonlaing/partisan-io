package main

import (
	"log"
	"net/http"
	"os"
	apiV1 "partisan/api/v1"
	apiV2 "partisan/api/v2"
	"partisan/auth"
	"partisan/db"
	"partisan/logger"
	m "partisan/models" // V1 models
	"time"

	"partisan/models.v2/users"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func init() {
	apiV1.ConfigureEmailer(emailConfig)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte("aoisahdfasodsaoih1289y3sopa0912"))
	r.Use(sessions.Sessions("partisan-io", store))
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

	v1Root := "api/v1"

	// V1
	{

		// allow users to pre sign up before the app is released
		r.POST(v1Root+"/presignup", apiV1.PreSignUpCreate)

		r.POST(v1Root+"/login", apiV1.LoginHandler)
		r.DELETE(v1Root+"/logout", apiV1.LogoutHandler)
		r.GET(v1Root+"/logout", apiV1.LogoutHandler)

		feed := r.Group(v1Root + "/feed")
		feed.Use(auth.Auth())
		{
			feed.GET("/", apiV1.FeedIndex)
			feed.GET("/socket", apiV1.FeedSocket)
			feed.GET("/show/:user_id", apiV1.FeedShow)
		}

		// DEPRECATED!!!!!!!!!!!
		users := r.Group(v1Root + "/users")
		users.Use(deprecated(v2Root + "/users"))
		users.Use(auth.Auth())
		{
			r.POST(v1Root+"/users", apiV1.UserCreate)
			r.GET(v1Root+"/user/check_unique", apiV1.UserCheckUnique)
			r.GET(v1Root+"/username_suggest", auth.Auth(), apiV1.UsernameSuggest)
			users.GET("/", apiV1.UserShow) // Show Current User
			users.PATCH("/", apiV1.UserUpdate)
			users.GET("/:user_id/match", apiV1.UserMatch)
			users.POST("/avatar_upload", apiV1.UserAvatarUpload)
		}

		profiles := r.Group(v1Root + "/profiles")
		profiles.Use(auth.Auth())
		{
			profiles.GET("/", apiV1.ProfileShow)         // Show Current User's profile
			profiles.GET("/:user_id", apiV1.ProfileShow) // Show Other User's profile
		}

		profile := r.Group(v1Root + "/profile")
		profile.Use(auth.Auth())
		{
			profile.PATCH("/", apiV1.ProfileUpdate) // Update Current User's profile
		}

		friends := r.Group(v1Root + "/friendships")
		friends.Use(auth.Auth())
		{
			friends.GET("/", apiV1.FriendshipIndex)
			friends.POST("/", apiV1.FriendshipCreate)
			friends.GET("/:friend_id", apiV1.FriendshipShow)
			friends.PATCH("/", apiV1.FriendshipConfirm)
			friends.DELETE("/", apiV1.FriendshipDestroy)
		}

		questions := r.Group(v1Root + "/questions")
		questions.Use(auth.Auth())
		{
			questions.GET("/", apiV1.QuestionIndex)
		}

		answers := r.Group(v1Root + "/answers")
		answers.Use(auth.Auth())
		{
			answers.PATCH("/", apiV1.AnswersUpdate)
		}

		posts := r.Group(v1Root + "/posts")
		posts.Use(auth.Auth())
		{
			// posts.GET("/", apiV1.PostsIndex)
			posts.POST("/", apiV1.PostsCreate)
			posts.GET("/:record_id", apiV1.PostsShow)
			posts.PATCH("/:id", apiV1.PostsUpdate)
			posts.DELETE("/:id", apiV1.PostsDestroy)
			posts.GET("/:record_id/likes", apiV1.LikeCount)
			posts.POST("/:record_id/likes", apiV1.LikeCreate)

			posts.GET("/:record_id/comments", apiV1.CommentsIndex)
			posts.GET("/:record_id/comments/count", apiV1.CommentsCount)

			posts.GET("/:record_id/attachments", apiV1.ImageAttachmentIndex)
			// posts.POST("/:record_id/attachments", apiV1.ImageAttachmentCreate)
		}

		comments := r.Group(v1Root + "/comments")
		comments.Use(auth.Auth())
		{
			comments.POST("/", apiV1.CommentsCreate)
			comments.GET("/:record_id/likes", apiV1.LikeCount)
			comments.POST("/:record_id/likes", apiV1.LikeCreate)
		}

		matches := r.Group(v1Root + "/matches")
		matches.Use(auth.Auth())
		{
			matches.GET("/", apiV1.MatchesIndex)
		}

		notifications := r.Group(v1Root + "/notifications")
		notifications.Use(auth.Auth())
		{
			notifications.GET("/", apiV1.NotificationsIndex)
			notifications.PATCH("/", apiV1.NotificationsRead)
			notifications.GET("/count", apiV1.NotificationsCount)
		}

		messages := r.Group(v1Root + "/messages")
		messages.Use(auth.Auth())
		{
			messages.GET("/threads", apiV1.MessageThreadIndex)
			messages.POST("/threads", apiV1.MessageThreadCreate)
			messages.GET("/count", apiV1.MessageCount)
			messages.GET("/threads/:thread_id", apiV1.MessageIndex)
			messages.POST("/threads/:thread_id", apiV1.MessageCreate)
			messages.GET("/threads/:thread_id/socket", apiV1.MessageSocket)
		}

		r.GET(v1Root+"/socket_ticket", auth.Auth(), apiV1.SocketTicketCreate)

		r.GET(v1Root+"/hashtags", auth.Auth(), apiV1.HashtagShow)

		r.POST(v1Root+"/flag", auth.Auth(), apiV1.FlagCreate)

	}

	// V2
	v2Root := "api/v2"
	{
		users := r.Group(v2Root + "/users")
		users.Use(auth.Auth())
		{
			r.POST(v2Root+"/users", apiV2.UserCreate)
			r.GET(v2Root+"/username_suggest", auth.Auth(), apiV2.UsernameSuggest)
			users.GET("/", apiV2.UserShow) // Show Current User
			users.PATCH("/", apiV2.UserUpdate)
			users.POST("/avatar_upload", apiV2.UserAvatarUpload)
		}

		// posts := r.Group(v2Root + "/posts")
		// posts.Use(auth.Auth())
		// {
		// 	posts.GET("/", apiV2.PostsIndex)
		// 	posts.POST("/", apiV2.PostsCreate)
		// 	posts.GET("/:record_id", apiV2.PostsShow)
		// 	posts.PATCH("/:id", apiV2.PostsUpdate)
		// 	posts.DELETE("/:id", apiV2.PostsDestroy)
		// 	posts.GET("/:record_id/likes", apiV2.LikeCount)
		// 	posts.POST("/:record_id/likes", apiV2.LikeCreate)

		// 	posts.GET("/:record_id/comments", apiV2.CommentsIndex)
		// 	posts.GET("/:record_id/comments/count", apiV2.CommentsCount)

		// 	posts.GET("/:record_id/attachments", apiV2.ImageAttachmentIndex)
		// 	// posts.POST("/:record_id/attachments", apiV1.ImageAttachmentCreate)
		// }
	}

	// HTML
	r.HTMLRender = createMyRender()

	r.GET("/profiles/:username", auth.Auth(), ProfileShow)
	r.GET("/feed", auth.Auth(), FeedIndex)
	r.GET("/profile", auth.Auth(), ProfileEdit)
	r.GET("/questions", auth.Auth(), QuestionsIndex)
	r.GET("/matches", auth.Auth(), MatchesIndex)
	r.GET("/friends", auth.Auth(), FriendsIndex)
	r.GET("/messages", auth.Auth(), MessagesIndex)
	r.GET("/comments/:record_id", auth.Auth(), CommentShow)
	r.GET("/likes/:record_id", auth.Auth(), LikeShow)
	r.GET("/posts/:record_id", auth.Auth(), PostShow)

	r.GET("/hashtags", auth.Auth(), HashtagShow)

	// no login on website anymore, mobile-only
	// r.GET("/login", Login)
	// r.GET("/signup", SignUp)

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
		&m.Post{},
		&users.User{},
		&m.Friendship{},
		&m.FeedItem{},
		&m.Like{},
		&m.Profile{},
		&m.Comment{},
		&m.ImageAttachment{},
		&m.Notification{},
		&m.Hashtag{},
		&m.Taxonomy{},
		&m.Flag{},
		&m.UserTag{},
		&m.Message{},
		&m.MessageThread{},
		&m.MessageThreadUser{},
		&m.SocketTicket{},
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

func createMyRender() multitemplate.Render {
	root := "templates"
	base := root + "/layout.html"

	r := multitemplate.New()
	r.AddFromFiles("feed", base, root+"/feed.html")
	r.AddFromFiles("hashtags", base, root+"/hashtags.html")
	r.AddFromFiles("login", base, root+"/login.html")
	r.AddFromFiles("matches", base, root+"/matches.html")
	r.AddFromFiles("messages", base, root+"/messages.html")
	r.AddFromFiles("friends", base, root+"/friends.html")
	r.AddFromFiles("post", base, root+"/post.html")
	r.AddFromFiles("profile_edit", base, root+"/profile_edit.html")
	r.AddFromFiles("profile_show", base, root+"/profile_show.html")
	r.AddFromFiles("questions", base, root+"/questions.html")
	r.AddFromFiles("signup", base, root+"/sign-up.html")

	return r
}

func deprecated(newPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if newPath != "" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorln("This endpoint is deprecated. Please see:", newPath))
			return
		}

		c.AbortWithError(http.StatusBadRequest, fmt.Errorln("This endpoint is deprecated."))
	}
}
