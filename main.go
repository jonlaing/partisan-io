package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	api "partisan/api/v1"
	"partisan/auth"
	"partisan/db"
	m "partisan/models"
)

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte("aoisahdfasodsaoih1289y3sopa0912"))
	r.Use(sessions.Sessions("partisan-io", store))

	// V1
	{
		v1Root := "api/v1"

		r.POST(v1Root+"/login", api.LoginHandler)
		r.DELETE(v1Root+"/logout", api.LogoutHandler)
		r.GET(v1Root+"/logout", api.LogoutHandler)

		feed := r.Group(v1Root + "/feed")
		feed.Use(auth.Auth())
		{
			feed.GET("/", api.FeedIndex)
		}

		users := r.Group(v1Root + "/users")
		users.Use(auth.Auth())
		{
			r.POST(v1Root+"/users", api.UserCreate)
			users.GET("/", api.UserShow) // Show Current User
			users.GET("/:user_id/match", api.UserMatch)
		}

		profiles := r.Group(v1Root + "/profiles")
		profiles.Use(auth.Auth())
		{
			profiles.GET("/", api.ProfileShow)         // Show Current User's profile
			profiles.GET("/:user_id", api.ProfileShow) // Show Other User's profile
			// profiles.PATCH("/:user_id", ProfileUpdate) // Show Other User's profile
		}

		friends := r.Group(v1Root + "/friends")
		friends.Use(auth.Auth())
		{
			friends.POST("/", api.FriendshipCreate)
			friends.POST("/confirm", api.FriendshipConfirm)
			friends.DELETE("/", api.FriendshipDestroy)
		}

		questions := r.Group(v1Root + "/questions")
		questions.Use(auth.Auth())
		{
			questions.GET("/", QuestionShow)
		}

		answers := r.Group(v1Root + "/answers")
		answers.Use(auth.Auth())
		{
			answers.PATCH("/", api.AnswersUpdate)
		}

		posts := r.Group(v1Root + "/posts")
		posts.Use(auth.Auth())
		{
			posts.GET("/", api.PostsIndex)
			posts.POST("/", api.PostsCreate)
			// posts.GET("/:id", api.PostsShow)
			posts.PATCH("/:id", api.PostsUpdate)
			posts.DELETE("/:id", api.PostsDestroy)
			posts.GET("/:record_id/likes", api.LikeCount)
			posts.POST("/:record_id/likes", api.LikeCreate)

			posts.GET("/:record_id/comments", api.CommentsIndex)
			posts.GET("/:record_id/comments/count", api.CommentsCount)
		}

		comments := r.Group(v1Root + "/comments")
		comments.Use(auth.Auth())
		{
			comments.POST("/", api.CommentsCreate)
			comments.GET("/:record_id/likes", api.LikeCount)
			comments.POST("/:record_id/likes", api.LikeCreate)
		}

	}

	// HTML
	r.LoadHTMLGlob("templates/*")
	htmlProfiles := r.Group("/profiles")
	htmlProfiles.Use(auth.Auth())
	{
		// htmlProfiles.GET("/", ProfileHTMLShowCurrent) // Will show editing options
		htmlProfiles.GET("/:user_id", ProfileShow)
	}

	r.Use(static.Serve("/", static.LocalFile("dist", false)))

	// Database Operations
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&m.Post{}, &m.User{}, &m.Friendship{}, &m.FeedItem{}, &m.Like{}, &m.Profile{}, &m.Comment{})

	r.Run(":4000")
}
