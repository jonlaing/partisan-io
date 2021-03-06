package main

import (
	apiV2 "partisan/api/v2"
	"partisan/auth"

	"github.com/gin-gonic/gin"
)

func init() {
	apiV2.ConfigureEmailer(emailConfig)
}

func initRoutesV2(r *gin.Engine) {
	// V2
	v2 := r.Group("api/v2")
	v2.Use(auth.AppToken())
	{
		// AUTH
		v2.POST("/login", apiV2.LoginHandler)
		v2.DELETE("/logout", apiV2.LogoutHandler)

		// USERS
		users := v2.Group("/users")
		users.Use(auth.Auth())
		{
			v2.POST("/users", apiV2.UserCreate)
			// v2.GET("/username_suggest", auth.Auth(), apiV2.UsernameSuggest)
			users.GET("/", apiV2.UserShow)         // Show Current User
			users.GET("/:user_id", apiV2.UserShow) // Show Other User
			users.PATCH("/", apiV2.UserUpdate)
			users.POST("/avatar_upload", apiV2.UserAvatarUpload)
		}

		// FRIENDSHIPS
		friends := v2.Group("/friendships")
		friends.Use(auth.Auth())
		{
			friends.GET("/", apiV2.FriendshipIndex)
			friends.POST("/", apiV2.FriendshipCreate)
			friends.GET("/:friend_id", apiV2.FriendshipShow)
			friends.PATCH("/:friend_id", apiV2.FriendshipUpdate)
			friends.DELETE("/", apiV2.FriendshipDestroy)
		}

		// POSTS
		posts := v2.Group("/posts")
		posts.Use(auth.Auth())
		{
			posts.GET("/", apiV2.PostIndex)
			posts.POST("/", apiV2.PostCreate)
			posts.GET("/:record_id", apiV2.PostShow)
			posts.PATCH("/:record_id", apiV2.PostUpdate)
			posts.DELETE("/:record_id", apiV2.PostDestroy)
			posts.POST("/:record_id/like", apiV2.LikeCreate)

			posts.GET("/:record_id/comments", apiV2.CommentIndex)
			posts.POST("/:record_id/comments", apiV2.CommentCreate)
		}

		// NOTIFICATIONS
		notifications := v2.Group("/notifications")
		notifications.Use(auth.Auth())
		{
			notifications.GET("/", apiV2.NotificationIndex)
			notifications.GET("/count", apiV2.NotificationsCount)
		}

		// MESSAGES
		messages := v2.Group("/messages")
		messages.Use(auth.Auth())
		{
			messages.GET("/threads", apiV2.ThreadIndex)
			messages.POST("/threads", apiV2.ThreadCreate)
			messages.GET("/threads/:thread_id", apiV2.MessageIndex)
			messages.POST("/threads/:thread_id", apiV2.MessageCreate)
			messages.GET("/threads/:thread_id/subscribe", apiV2.MessageThreadSubscribe)
			messages.GET("/threads/:thread_id/new", apiV2.NewMessages)
			messages.GET("/unread", apiV2.MessageUnread)
		}

		// EVENTS
		events := v2.Group("/events")
		events.Use(auth.Auth())
		{
			events.GET("/", apiV2.EventIndex)
			events.POST("/", apiV2.EventCreate)
			events.GET("/:event_id", apiV2.EventShow)
			events.PATCH("/:event_id", apiV2.EventUpdate)
			events.POST("/:event_id/hosts/:user_id", apiV2.EventAddHost)
			events.DELETE("/:event_id/hosts/:user_id", apiV2.EventRemoveHost)
			events.POST("/:event_id/going", apiV2.EventGoing)
			events.POST("/:event_id/maybe", apiV2.EventMaybe)
			events.DELETE("/:event_id/unsubscribe", apiV2.EventUnsubscribe)
			events.DELETE("/:event_id", apiV2.EventDestroy)

			events.GET("/:event_id/posts", apiV2.EventPosts)
		}

		// COMMENTS
		v2.POST("/comments/:record_id/like", auth.Auth(), apiV2.LikeCreate)

		// MATCHES
		v2.POST("/matches", auth.Auth(), apiV2.MatchIndex)

		// ANSWERS
		v2.PATCH("/answers", auth.Auth(), apiV2.AnswersUpdate)

		// QUESTIONS
		v2.GET("/questions", auth.Auth(), apiV2.QuestionIndex)

		// FLAGS
		v2.POST("/flag", auth.Auth(), apiV2.FlagCreate)

		// HASHTAGS
		v2.GET("/search", auth.Auth(), apiV2.HashtagShow)

		v2.GET("/socket_ticket", auth.Auth(), apiV2.SocketTicketCreate)

		// should only be able to request a password reset to the app
		v2.POST("/password_reset", apiV2.PasswordResetCreate)
	}

	// by virtue of not being logged in, and currently having to do this from the website, we can't use
	// the app token here
	r.GET("api/v2/password_reset/:reset_id", apiV2.PasswordResetShow)
	r.PATCH("api/v2/password_reset/:reset_id", apiV2.PasswordResetUpdate)

	// view a post without being logged in
	// NOTE: Don't show user details here!
	r.GET("posts/:post_id", PostShow)
	r.HEAD("posts/:post_id", PostShow)
}
