package main

import (
	"net/http"
	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// Login shows the login screen
func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{"title": "Login to Partisan.IO"})
}

// SignUp shows the signup screen
func SignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "signup", gin.H{"title": "Sign Up for Partisan.IO"})
}
