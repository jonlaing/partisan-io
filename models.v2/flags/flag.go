package flags

import (
	"time"
)

const (
	FlNoReason  string = ""
	FlOffensive        = "offensive"
	FlCopyright        = "copyright"
	FlSpam             = "spam"
	FlOther            = "other"
)

// Flag alerts admins of potentially problematic behavior of a user,
// or problematic content in a post/comment
type Flag struct {
	ID         string    `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	UserID     string    `json:"user_id" sql:"type:uuid;default:uuid_generate_v4()"`
	RecordID   string    `json:"record_id" sql:"type:uuid;default:uuid_generate_v4()"`
	RecordType string    `json:"record_type"`
	Reason     string    `json:"reason"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
