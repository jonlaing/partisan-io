package friendships

import "github.com/jinzhu/gorm"

func GetByUserID(userID string, db *gorm.DB) (fs []Friendship, err error) {
	err = db.Where("user_id = ? OR friend_id = ?", userID, userID).Find(&fs).Error
	return
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

func GetConfirmedByUserID(userID string, db *gorm.DB) (fs []Friendship, err error) {
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
