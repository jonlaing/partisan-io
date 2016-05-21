package friendships

import "github.com/jinzhu/gorm"

func ListByUserID(userID string, db *gorm.DB) (fs Friendships, err error) {
	err = db.Where("user_id = ? OR friend_id = ?", userID, userID).Find(&fs).Error
	// don't collect the users here because we do that in the handler
	return
}

func GetByUserIDs(userID, friendID string, db *gorm.DB) (f Friendship, err error) {
	if err = db.Where("user_id = ? AND friend_id = ?", userID, friendID).Find(&f).Error; err != nil {
		// if you don't find it the first way, try the other way
		if err = db.Where("user_id = ? AND friend_id = ?", friendID, userID).Find(&f).Error; err != nil {
			return
		}
	}

	return
}

func Exists(userID, friendID string, db *gorm.DB) bool {
	var count int
	if err := db.Table("friendships").
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Count(&count).Error; err == nil && count > 0 {
		return true
	}

	if err := db.Table("friendships").
		Where("user_id = ? AND friend_id = ?", friendID, userID).
		Count(&count).Error; err == nil && count > 0 {
		return true
	}

	return false
}

func GetIDsByUserID(userID string, db *gorm.DB) (uuids []string, err error) {
	fids := []string{}
	uids := []string{}

	err = db.Table("friendships").Where("user_id = ?", userID).Pluck("friend_id", &fids).Error
	if err != nil {
		return
	}

	err = db.Table("friendships").Where("friend_id = ?", userID).Pluck("user_id", &uids).Error
	if err != nil {
		return
	}

	uuids = removeDuplicates(append(fids, uids...))
	return
}

func GetConfirmedByUserID(userID string, db *gorm.DB) (fs Friendships, err error) {
	err = db.Where("user_id = ? OR friend_id = ?", userID, userID).Where("confirmed = ?", true).Find(&fs).Error

	return
}

func GetConfirmedIDsByUserID(userID string, db *gorm.DB) (uuids []string, err error) {
	fids := []string{}
	uids := []string{}

	err = db.Table("friendships").Where("user_id = ?", userID).Where("confirmed = ?", true).Pluck("friend_id", &fids).Error
	if err != nil {
		return
	}

	err = db.Table("friendships").Where("friend_id = ?", userID).Where("confirmed = ?", true).Pluck("user_id", &uids).Error
	if err != nil {
		return
	}

	uuids = removeDuplicates(append(fids, uids...))
	return
}

func removeDuplicates(in []string) (out []string) {
	found := make(map[string]bool)
	for _, s := range in {
		if _, ok := found[s]; !ok {
			found[s] = true
			out = append(out, s)
		}
	}

	return
}
