package main

import (
	"net/http"

	"gin-gonic/contrib/sessions"
	"gin-gonic/gin"
)

// Login shows the login screen
func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{"title": "Login to Partisan.IO"})
}

// SignUp shows the signup screen
func SignUp(c *gin.Context) {
	sess := sessions.Default(c)

	if sess.Get("user_id") != nil {
		c.Redirect(http.StatusFound, "/feed/")
		return
	}

	c.HTML(http.StatusOK, "signup", gin.H{"title": "Sign Up for Partisan.IO"})
}
