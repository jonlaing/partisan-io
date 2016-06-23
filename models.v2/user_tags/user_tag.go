package usertags

import (
	"regexp"

	"partisan/logger"

	"partisan/models.v2/friendships"
	"partisan/models.v2/notifications"
	"partisan/models.v2/posts"
	"partisan/models.v2/users"

	"github.com/jinzhu/gorm"
)

// Usertagger is an interface for records that can be tagged with a username
type Usertagger interface {
	GetID() string
	GetUserID() string
	GetType() string
	GetContent() string
}

type usertag struct {
	recordID string
}

func (t usertag) GetID() string {
	return t.recordID
}

func (t usertag) GetAction() string {
	return string(notifications.AUserTag)
}

func Extract(r Usertagger, db *gorm.DB) error {
	tags := extractTags(r.GetContent())
	if len(tags) == 0 {
		return ErrNoTags
	}

	us, err := users.ListByUsernames(db, tags...)
	if err != nil {
		return err
	}

	for _, user := range us {
		// can't tag users that aren't your friend
		if friendships.Exists(r.GetUserID(), user.ID, db) {
			createPost(user.ID, r, db)
		}
	}

	return nil
}

// create a post with action AUserTag, so user tags will show up in feed,
// and then notify the user that was tagged
func createPost(userID string, r Usertagger, db *gorm.DB) {
	binding := posts.CreatorBinding{
		ParentType: r.GetType(),
		ParentID:   r.GetID(),
		Action:     posts.AUserTag,
	}
	p, errs := posts.New(userID, binding)
	if len(errs) == 0 {
		if err := db.Create(&p).Error; err == nil {
			notifyUser(userID, r, db)
		}
	} else {
		logger.Error.Println(errs)
	}
}

func notifyUser(userID string, r Usertagger, db *gorm.DB) {
	n, errs := notifications.New(r.GetUserID(), userID, usertag{r.GetID()})
	if len(errs) == 0 {
		db.Create(&n)
	}
}

// extractTags finds usertags in a string
func extractTags(text string) []string {
	usertagSearch := regexp.MustCompile("@([a-zA-Z0-9_.]+)")
	excludeEmail := regexp.MustCompile("\\.[a-zA-Z]+")
	usertagSlices := usertagSearch.FindAllStringSubmatch(text, -1) // -1 means find all instances

	var usertags []string

	for _, usertag := range usertagSlices {
		if !excludeEmail.MatchString(usertag[1]) {
			usertags = append(usertags, usertag[1])
		}
	}

	return usertags
}
