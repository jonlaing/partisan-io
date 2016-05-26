package main

import (
	apiV2 "partisan/api/v2"
	"partisan/auth"

	"github.com/gin-gonic/gin"
)

func initRoutesV1(r *gin.Engine) {
	// V2
	v2Root := "api/v2"
	{
		// AUTH
		r.POST(v2Root+"/login", apiV2.LoginHandler)
		r.DELETE(v2Root+"/logout", apiV2.LogoutHandler)

		// USERS
		users := r.Group(v2Root + "/users")
		users.Use(auth.Auth())
		{
			r.POST(v2Root+"/users", apiV2.UserCreate)
			r.GET(v2Root+"/username_suggest", auth.Auth(), apiV2.UsernameSuggest)
			users.GET("/", apiV2.UserShow)          // Show Current User
			users.GET("/:username", apiV2.UserShow) // Show Other User
			users.PATCH("/", apiV2.UserUpdate)
			users.POST("/avatar_upload", apiV2.UserAvatarUpload)
		}

		// FRIENDSHIPS
		friends := r.Group(v2Root + "/friendships")
		friends.Use(auth.Auth())
		{
			friends.GET("/", apiV2.FriendshipIndex)
			friends.POST("/", apiV2.FriendshipCreate)
			friends.GET("/:friend_id", apiV2.FriendshipShow)
			friends.PATCH("/", apiV2.FriendshipConfirm)
			friends.DELETE("/", apiV2.FriendshipDestroy)
		}

		// POSTS
		posts := r.Group(v2Root + "/posts")
		posts.Use(auth.Auth())
		{
			posts.GET("/", apiV2.PostsIndex)
			posts.POST("/", apiV2.PostsCreate)
			posts.GET("/:record_id", apiV2.PostsShow)
			posts.PATCH("/:record_id", apiV2.PostsUpdate)
			posts.DELETE("/:record_id", apiV2.PostsDestroy)
			posts.POST("/:record_id/like", apiV2.LikeCreate)

			posts.GET("/:record_id/comments", apiV2.CommentsIndex)
			posts.POST("/:record_id/comments", apiV2.CommentsCreate)

			posts.GET("/:record_id/attachments", apiV2.ImageAttachmentIndex)
			// posts.POST("/:record_id/attachments", apiV2.ImageAttachmentCreate)
		}

		// COMMENTS
		r.POST(v2Root+"/comments/:record_id/like", auth.Auth(), apiV2.LikeCreate)

		// MATCHES
		r.GET(v2Root+"/matches", auth.Auth(), apiV2.MatchesIndex)

		// ANSWERS
		r.PATCH(v2Root+"/answers", auth.Auth(), apiV2.AnswersUpdate)

		// QUESTIONS
		r.GET(v2Root+"/questions", auth.Auth(), apiV2.QuestionIndex)

		// FLAGS
		r.POST(v2Root+"/flag", auth.Auth(), apiV2.FlagCreate)

		// HASHTAGS
		r.GET(v1Root+"/search", auth.Auth(), apiV1.HashtagShow)
	}
}
