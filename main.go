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
			questions.GET("/", QuestionsTest)
			// questions.GET("/", QuestionsIndex)
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

	r.Use(static.Serve("/", static.LocalFile("dist", false)))

	// Database Operations
	db, err := initDB()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Post{}, &User{}, &Friendship{}, &FeedItem{}, &Like{}, &Dislike{})

	r.Run(":4000")
}
