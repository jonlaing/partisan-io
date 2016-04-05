package emailer

import (
	"fmt"
	"net/smtp"

	"github.com/scorredoira/email"
)

// Config is configuration for the emaail
type Config struct {
	From     string
	Password string
}

var config Config

// Configure copies the a configuration to the package
func Configure(c Config) {
	config = c
}

// SendEmail sends an email through gmail
func SendEmail(to, subject, body string) error {
	m := email.NewMessage(subject, body)
	m.From = config.From
	m.To = []string{to}

	fmt.Println(string(m.Bytes()))
	return email.Send("smtp.gmail.com:587", smtp.PlainAuth("", config.From, config.Password, "smtp.gmail.com"), m)
}

// SendWelcomeEmail sends the welcome email to a user
func SendWelcomeEmail(username, email string) error {
	subject := "Welcome to Partisan.IO"
	bodyTemp := `
	Welcome to Partisan.IO

	Your account has been created! Your username is: @%s

	Have fun!

	- Jon at Partisan.IO`

	body := fmt.Sprintf(bodyTemp, username)

	return SendEmail(email, subject, body)
}
