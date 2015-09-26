package v1

import (
	"partisan/emailer"
)

// ConfigureEmailer sets the emailer configuration
func ConfigureEmailer(config emailer.Config) {
	emailer.Configure(config)
}
