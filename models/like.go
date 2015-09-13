package models

// Like is polymorphic
type Like struct {
	ID         uint64 `gorm:"primary_key"`
	UserID     uint64
	RecordID   uint64
	RecordType string
	IsDislike  bool // so that we can use the same table for both likes and dislikes, not currently in use
}
