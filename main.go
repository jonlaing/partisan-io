package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte("aoisahdfasodsaoih1289y3sopa0912"))
	r.Use(sessions.Sessions("partisan-io", store))

	// V1
	{
		v1Root := "api/v1"

		r.POST(v1Root+"/login", LoginHandler)
		r.DELETE(v1Root+"/logout", LogoutHandler)
		r.GET(v1Root+"/logout", LogoutHandler)

		feed := r.Group(v1Root + "/feed")
		feed.Use(Auth())
		{
			feed.GET("/", FeedIndex)
		}

		users := r.Group(v1Root + "/users")
		users.Use(Auth())
		{
			r.POST(v1Root+"/users", UserCreate)
			users.GET("/", UserShow) // Show Current User
			users.GET("/:user_id/match", UserMatch)
		}

		profiles := r.Group(v1Root + "/profiles")
		profiles.Use(Auth())
		{
			profiles.GET("/", ProfileShow)         // Show Current User's profile
			profiles.GET("/:user_id", ProfileShow) // Show Other User's profile
			// profiles.PATCH("/:user_id", ProfileUpdate) // Show Other User's profile
		}

		friends := r.Group(v1Root + "/friends")
		friends.Use(Auth())
		{
			friends.POST("/", FriendshipCreate)
			friends.POST("/confirm", FriendshipConfirm)
			friends.DELETE("/", FriendshipDestroy)
		}

		questions := r.Group(v1Root + "/questions")
		questions.Use(Auth())
		{
			questions.GET("/", QuestionShow)
		}

		answers := r.Group(v1Root + "/answers")
		answers.Use(Auth())
		{
			answers.PATCH("/", AnswersUpdate)
		}

		posts := r.Group(v1Root + "/posts")
		posts.Use(Auth())
		{
			posts.GET("/", PostsIndex)
			posts.POST("/", PostsCreate)
			posts.GET("/:id", PostsShow)
			posts.PATCH("/:id", PostsUpdate)
			posts.DELETE("/:id", PostsDestroy)
			posts.POST("/:post_id/like", LikeCreate)
			posts.POST("/:post_id/dislike", DislikeCreate)
		}

	}

	// HTML
	r.LoadHTMLGlob("templates/*")
	htmlProfiles := r.Group("/profiles")
	htmlProfiles.Use(Auth())
	{
		// htmlProfiles.GET("/", ProfileHTMLShowCurrent) // Will show editing options
		htmlProfiles.GET("/:user_id", ProfileHTMLShow)
	}

	r.Use(static.Serve("/", static.LocalFile("dist", false)))

	// Database Operations
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Post{}, &User{}, &Friendship{}, &FeedItem{}, &Like{}, &Dislike{}, &Profile{})

	r.Run(":4000")
}
