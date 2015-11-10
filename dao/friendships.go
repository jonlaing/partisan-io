package dao

import (
	m "partisan/models"

	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"
)

// Friends returns all the User's friends
func Friends(u m.User, confirmed bool, db *gorm.DB) ([]m.User, error) {
	friendIDs, err := FriendIDs(u, confirmed, db)
	if err != nil {
		return []m.User{}, err
	}

	users := []m.User{}
	db.Where(friendIDs).Order("username asc").Find(&users)
	return users, nil
}

// ConfirmedFriends returns all the User's confirmed friends
func ConfirmedFriends(u m.User, db *gorm.DB) ([]m.User, error) {
	return Friends(u, true, db)
}

// FriendIDs returns all IDs of User's Friends
func FriendIDs(u m.User, confirmed bool, db *gorm.DB) ([]uint64, error) {
	var friendIDs []uint64
	var friendships []m.Friendship
	err := db.Where("(user_id = ? OR friend_id = ?) AND confirmed = ?", u.ID, u.ID, confirmed).Find(&friendships).Error
	if err != nil {
		return []uint64{}, err
	}

	for _, f := range friendships {
		if f.UserID != u.ID {
			friendIDs = append(friendIDs, f.UserID)
		}

		if f.FriendID != u.ID {
			friendIDs = append(friendIDs, f.FriendID)
		}
	}

	return friendIDs, err
}

// ConfirmedFriendIDs returns all IDs of User's confirmed friends
func ConfirmedFriendIDs(u m.User, db *gorm.DB) ([]uint64, error) {
	return FriendIDs(u, true, db)
}

func GetFriendship(u m.User, fID uint64, db *gorm.DB) (m.Friendship, error) {
	var f1, f2 m.Friendship

	if err := db.Where("user_id = ? AND friend_id = ?", u.ID, fID).First(&f1).Error; err == nil {
		return f1, nil
	}

	err := db.Where("user_id = ? AND friend_id = ?", fID, u.ID).First(&f2).Error

	return f2, err
}
