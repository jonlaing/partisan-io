package v2

import (
	"fmt"
	"partisan/emailer"

	"github.com/timehop/apns"
)

var pushClient apns.Client

func init() {
	var err error
	pushClient, err = apns.NewClientWithFiles(apns.SandboxGateway, "pushcert.pem", "pushkey.pem")
	if err != nil {
		fmt.Println("Couldn't connect to APNS:", err)
	}

	go func() {
		for f := range pushClient.FailedNotifs {
			fmt.Println("Notif", f.Notif.ID, "failed with", f.Err.Error())
			fmt.Println("Device Token:", f.Notif.DeviceToken)
		}
	}()
}

func ConfigureEmailer(config emailer.Config) {
	emailer.Configure(config)
}
