package main

import (
	"github.com/gin-gonic/contrib/renders/multitemplate"
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

	v1Root := "api/v1"

	// V1
	{

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
			users.PATCH("/", api.UserUpdate)
			users.GET("/:user_id/match", api.UserMatch)
			users.POST("/avatar_upload", api.UserAvatarUpload)
		}

		profiles := r.Group(v1Root + "/profiles")
		profiles.Use(auth.Auth())
		{
			profiles.GET("/", api.ProfileShow)         // Show Current User's profile
			profiles.GET("/:user_id", api.ProfileShow) // Show Other User's profile
		}

		profile := r.Group(v1Root + "/profile")
		profile.Use(auth.Auth())
		{
			profile.PATCH("/", api.ProfileUpdate) // Update Current User's profile
		}

		friends := r.Group(v1Root + "/friendships")
		friends.Use(auth.Auth())
		{
			friends.POST("/", api.FriendshipCreate)
			friends.GET("/:friend_id", api.FriendshipShow)
			friends.PATCH("/", api.FriendshipConfirm)
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

			posts.GET("/:record_id/attachments", api.ImageAttachmentIndex)
			// posts.POST("/:record_id/attachments", api.ImageAttachmentCreate)
		}

		comments := r.Group(v1Root + "/comments")
		comments.Use(auth.Auth())
		{
			comments.POST("/", api.CommentsCreate)
			comments.GET("/:record_id/likes", api.LikeCount)
			comments.POST("/:record_id/likes", api.LikeCreate)
		}

		matches := r.Group(v1Root + "/matches")
		matches.Use(auth.Auth())
		{
			matches.GET("/", api.MatchesIndex)
		}

		notifications := r.Group(v1Root + "/notifications")
		notifications.Use(auth.Auth())
		{
			notifications.GET("/", api.NotificationsIndex)
			notifications.PATCH("/", api.NotificationsRead)
			notifications.GET("/count", api.NotificationsCount)
		}

		r.GET(v1Root+"/hashtags", api.HashtagShow)

	}

	// HTML
	r.HTMLRender = createMyRender()

	r.GET("/profiles/:user_id", auth.Auth(), ProfileShow)
	r.GET("/feed", auth.Auth(), FeedIndex)
	r.GET("/profile", auth.Auth(), ProfileEdit)
	r.GET("/questions", auth.Auth(), QuestionsIndex)
	r.GET("/matches", auth.Auth(), MatchesIndex)
	r.GET("/comments/:record_id", auth.Auth(), CommentShow)
	r.GET("/likes/:record_id", auth.Auth(), LikeShow)
	r.GET("/posts/:record_id", auth.Auth(), PostShow)

	r.GET("/hashtags", auth.Auth(), HashtagShow)

	r.GET("/login", Login)
	r.GET("/signup", SignUp)

	r.Use(static.Serve("/", static.LocalFile("dist", false)))
	r.Use(static.Serve("/localfiles", static.LocalFile("localfiles", false)))

	// Database Operations
	db, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(
		&m.Post{},
		&m.User{},
		&m.Friendship{},
		&m.FeedItem{},
		&m.Like{},
		&m.Profile{},
		&m.Comment{},
		&m.ImageAttachment{},
		&m.Notification{},
		&m.Hashtag{},
		&m.Taxonomy{},
	)

	r.Run(":4000")
}

func createMyRender() multitemplate.Render {
	root := "templates"
	base := root + "/layout.html"

	r := multitemplate.New()
	r.AddFromFiles("feed", base, root+"/feed.html")
	r.AddFromFiles("hashtags", base, root+"/hashtags.html")
	r.AddFromFiles("login", base, root+"/login.html")
	r.AddFromFiles("matches", base, root+"/matches.html")
	r.AddFromFiles("post", base, root+"/post.html")
	r.AddFromFiles("profile_edit", base, root+"/profile_edit.html")
	r.AddFromFiles("profile_show", base, root+"/profile_show.html")
	r.AddFromFiles("questions", base, root+"/questions.html")
	r.AddFromFiles("signup", base, root+"/sign-up.html")

	return r
}
