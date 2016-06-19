package v2

import "partisan/emailer"

func ConfigureEmailer(config emailer.Config) {
	emailer.Configure(config)
}
