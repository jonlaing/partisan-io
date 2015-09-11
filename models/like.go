package models

// Like is polymorphic
type Like struct {
	ID         uint64 `gorm:"primary_key"`
	UserID     uint64
	RecordID   uint64
	RecordType string
}

// Dislike is polymorphic
type Dislike struct {
	ID         uint64 `gorm:"primary_key"`
	UserID     uint64
	RecordID   uint64
	RecordType string
}
