package main

import (
	apiV1 "partisan/api/v1"
	"partisan/auth"

	"github.com/gin-gonic/gin"
)

func initRoutesV1(r *gin.Engine) {
	v1 := r.Group("api/v1")
	// V1
	{
		// allow users to pre sign up before the app is released
		v1.POST("/presignup", apiV1.PreSignUpCreate)

		v2.POST("/login", deprecated("api/v2/login"), apiV1.LoginHandler)
		v2.DELETE("/logout", deprecated("api/v2/logout"), apiV1.LogoutHandler)
		v2.GET("/logout", deprecated("api/v2/logout"), apiV1.LogoutHandler)

		feed := v1.Group("/feed")
		feed.User(deprecated("api/v2/posts"))
		feed.Use(auth.Auth())
		{
			feed.GET("/", apiV1.FeedIndex)
			feed.GET("/socket", apiV1.FeedSocket)
			feed.GET("/show/:user_id", apiV1.FeedShow)
		}

		// DEPRECATED!!!!!!!!!!!
		users := v1.Group("/users")
		users.Use(deprecated("api/v2/users"))
		users.Use(auth.Auth())
		{
			v1.POST("/users", apiV1.UserCreate)
			v1.GET("/user/check_unique", apiV1.UserCheckUnique)
			v1.GET("/username_suggest", auth.Auth(), apiV1.UsernameSuggest)
			users.GET("/", apiV1.UserShow) // Show Current User
			users.PATCH("/", apiV1.UserUpdate)
			users.GET("/:user_id/match", apiV1.UserMatch)
			users.POST("/avatar_upload", apiV1.UserAvatarUpload)
		}

		profiles := v1.Group("/profiles")
		users.Use(deprecated("api/v2/users"))
		profiles.Use(auth.Auth())
		{
			profiles.GET("/", apiV1.ProfileShow)         // Show Current User's profile
			profiles.GET("/:user_id", apiV1.ProfileShow) // Show Other User's profile
		}

		profile := v1.Group("/profile")
		profile.Use(deprecated("api/v2/users"))
		profile.Use(auth.Auth())
		{
			profile.PATCH("/", apiV1.ProfileUpdate) // Update Current User's profile
		}

		friends := v1.Group("/friendships")
		friends.Use(deprecated("api/v2/friendships"))
		friends.Use(auth.Auth())
		{
			friends.GET("/", apiV1.FriendshipIndex)
			friends.POST("/", apiV1.FriendshipCreate)
			friends.GET("/:friend_id", apiV1.FriendshipShow)
			friends.PATCH("/", apiV1.FriendshipConfirm)
			friends.DELETE("/", apiV1.FriendshipDestroy)
		}

		questions := v1.Group("/questions")
		questions.Use(deprecated("api/v2/questions"))
		questions.Use(auth.Auth())
		{
			questions.GET("/", apiV1.QuestionIndex)
		}

		answers := v1.Group("/answers")
		answers.Use(deprecated("api/v2/answers"))
		answers.Use(auth.Auth())
		{
			answers.PATCH("/", apiV1.AnswersUpdate)
		}

		posts := v1.Group("/posts")
		posts.Use(deprecated("api/v2/posts"))
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

		comments := v1.Group("/comments")
		comments.Use(deprecated("api/v2/comments"))
		comments.Use(auth.Auth())
		{
			comments.POST("/", apiV1.CommentsCreate)
			comments.GET("/:record_id/likes", apiV1.LikeCount)
			comments.POST("/:record_id/likes", apiV1.LikeCreate)
		}

		matches := v1.Group("/matches")
		matches.Use(deprecated("api/v2/matches"))
		matches.Use(auth.Auth())
		{
			matches.GET("/", apiV1.MatchesIndex)
		}

		notifications := v1.Group("/notifications")
		notifications.Use(deprecated("api/v2/notifications"))
		notifications.Use(auth.Auth())
		{
			notifications.GET("/", apiV1.NotificationsIndex)
			notifications.PATCH("/", apiV1.NotificationsRead)
			notifications.GET("/count", apiV1.NotificationsCount)
		}

		messages := v1.Group("/messages")
		messages.Use(deprecated("api/v2/messages"))
		messages.Use(auth.Auth())
		{
			messages.GET("/threads", apiV1.MessageThreadIndex)
			messages.POST("/threads", apiV1.MessageThreadCreate)
			messages.GET("/count", apiV1.MessageCount)
			messages.GET("/threads/:thread_id", apiV1.MessageIndex)
			messages.POST("/threads/:thread_id", apiV1.MessageCreate)
			messages.GET("/threads/:thread_id/socket", apiV1.MessageSocket)
		}

		v1.GET("/socket_ticket", auth.Auth(), deprecated("api/v2/socket_ticket"), apiV1.SocketTicketCreate)

		v1.GET("/hashtags", auth.Auth(), deprecated("api/v2/search"), apiV1.HashtagShow)

		v1.POST("/flag", auth.Auth(), deprecated("api/v2/flag"), apiV1.FlagCreate)

	}
}
