package emailer

import (
	"fmt"
	"net/smtp"

	"github.com/scorredoira/email"
)

// Config is configuration for the emaail
type Config struct {
	Auth     string
	Password string
	From     string
}

var config Config

// Configure copies the a configuration to the package
func Configure(c Config) {
	config = c
}

// SendEmail sends an email through gmail
func SendPlainEmail(to, subject, body string) error {
	m := email.NewMessage(subject, body)
	return send(to, m)
}

func SendHTMLEmail(to, subject, body string) error {
	m := email.NewHTMLMessage(subject, body)
	return send(to, m)
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

	return SendPlainEmail(email, subject, body)
}

func SendPasswordReset(email, resetID string) error {
	subject := "Partisan.IO: Password Reset"
	bodyTemp := `
	<p>
		A request was made to reset your password. If you did not request to reset your password, 
		disregard this email.
	</p>

	<p>To reset your password follow click the following link on your mobile device: <a href="partisanio://password_reset?reset_id=%s">Click Here</a></p>

	<p>- Jon at Partisan.IO</p>
	`

	body := fmt.Sprintf(bodyTemp, resetID)

	return SendHTMLEmail(email, subject, body)
}

func send(to string, m *email.Message) error {
	m.To = []string{to}

	return email.Send("smtp.gmail.com:587", smtp.PlainAuth("", config.Auth, config.Password, "smtp.gmail.com"), m)
}
