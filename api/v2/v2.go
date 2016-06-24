package v2

import (
	"partisan/emailer"
	"partisan/logger"

	"github.com/timehop/apns"
)

var pushClient apns.Client

func init() {
	var err error
	pushClient, err = apns.NewClientWithFiles(apns.SandboxGateway, "pushcert.pem", "pushkey.pem")
	if err != nil {
		logger.Error.Println("Couldn't connect to APNS:", err)
	}
}

func ConfigureEmailer(config emailer.Config) {
	emailer.Configure(config)
}
