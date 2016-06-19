package main

import (
	"os"
	"partisan/emailer"
)

var emailConfig = emailer.Config{
	Auth:     os.Getenv("EMAIL_UN"),
	Password: os.Getenv("EMAIL_PASS"),
	From:     "no-reply@partisan.io",
}
