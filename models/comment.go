package models

import(
  "time"
)

// Comment is for blah
type Comment struct {
  ID         uint64 `gorm:"primary_key" json:"id"`
  RecordType string `json:"record_type"`
  RecordID   uint64 `json:"record_id"`
  UserID     uint64 `json:"user_id"`
  Body       string `json:"body"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}
