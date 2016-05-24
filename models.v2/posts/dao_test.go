package posts

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"partisan/models.v2/flags"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"
)

var testdb *gorm.DB
var id string
var userID string
var p, comment, like, flagged Post

func init() {
	var err error
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PW"))
	testdb, err = gorm.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	if err := testdb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		panic(err)
	}

	testdb.AutoMigrate(Post{}, flags.Flag{})
}

func TestMain(m *testing.M) {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	userID = uuid.String()

	p = Post{
		UserID:    userID,
		Action:    APost,
		Body:      "my body",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := testdb.Save(&p).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&p)

	id = p.ID

	comment = Post{
		UserID:     userID,
		ParentID:   sql.NullString{id, true},
		ParentType: PTPost,
		Action:     AComment,
		Body:       "comment",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := testdb.Save(&comment).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&comment)

	like = Post{
		UserID:     userID,
		ParentID:   sql.NullString{id, true},
		ParentType: PTPost,
		Action:     ALike,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := testdb.Save(&like).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&like)

	flagged = Post{
		UserID:    userID,
		Action:    APost,
		Body:      "flagged post",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := testdb.Save(&flagged).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&flagged)

	flag := flags.Flag{
		UserID:     userID,
		RecordID:   flagged.ID,
		RecordType: "post",
	}

	if err := testdb.Save(&flag).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&flag)

	m.Run()
}

func TestGetByID(t *testing.T) {
	if id == "" {
		t.Error("Expected ID to have value")
		return
	}

	post, err := GetByID(id, userID, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
		return
	}

	if post.ID != id {
		t.Error("Expected post id to be:", id, ", got:", post.ID)
	}
}

func TestGetFeedByUserIDs(t *testing.T) {
	otherUserID, _ := uuid.NewV4()

	posts, err := GetFeedByUserIDs(userID, []string{otherUserID.String()}, 0, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(posts) > 0 {
		t.Error("Expected not to find posts")
	}

	posts, err = GetFeedByUserIDs(userID, []string{userID, otherUserID.String()}, 0, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(posts) != 1 {
		t.Error("Expected one post, got:", len(posts))
	}
}

func TestGetFeedByUserIDsAfter(t *testing.T) {
	otherUserID, _ := uuid.NewV4()

	posts, err := GetFeedByUserIDsAfter(userID, []string{userID, otherUserID.String()}, time.Now().Add(time.Hour), testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(posts) > 0 {
		t.Error("Expected not to find posts")
	}

	p.CreatedAt = time.Now().Add(time.Hour)
	p.UpdatedAt = p.CreatedAt

	if err := testdb.Save(&p).Error; err != nil {
		t.Error("Unexpected error saving post")
	}

	posts, err = GetFeedByUserIDsAfter(userID, []string{userID, otherUserID.String()}, time.Now(), testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(posts) != 1 {
		t.Error("Expected one post, got:", len(posts))
	}
}

func TestGetParent(t *testing.T) {
	parent, err := comment.GetParent(testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if parent.ID != p.ID {
		t.Error("Expected uuid to be:", p.ID, "got :", parent.ID)
	}
}

func TestGetChildren(t *testing.T) {
	children, err := p.GetChildren(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if len(children) != 2 {
		t.Error("Expected 2 children, got:", len(children))
	}

	foundComment := false
	foundLike := false
	for _, c := range children {
		if c.ID == comment.ID {
			foundComment = true
		}

		if c.ID == like.ID {
			foundLike = true
		}
	}

	if !foundComment {
		t.Error("Expected to find comment:", comment.ID, "in list")
	}

	if !foundLike {
		t.Error("Expected to find like:", like.ID, "in list")
	}
}

func TestGetLikeByUserID(t *testing.T) {
	l, err := p.GetLikeByUserID(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if l.ID != like.ID {
		t.Error("Expected likes to match")
	}
}

func TestGetLikeCount(t *testing.T) {
	err := p.GetLikeCount(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if p.LikeCount != 1 {
		t.Error("Expected like count to be 1, got:", p.LikeCount)
	}

	if !p.Liked {
		t.Error("Expected post to be liked")
	}

	err = comment.GetLikeCount(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if comment.LikeCount != 0 {
		t.Error("Expected like count to be 0, got:", comment.LikeCount)
	}

	if comment.Liked {
		t.Error("Expected comment not to be liked")
	}
}

func TestGetCommentCount(t *testing.T) {
	err := p.GetCommentCount(testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if p.CommentCount != 1 {
		t.Error("Expected comment count to be 1, got:", p.CommentCount)
	}
}

func TestGetChildCountList(t *testing.T) {
	ps := Posts{p}
	err := ps.GetChildCountList(userID, testdb)
	if err != nil {
		t.Error("Unexpected error")
	}

	if ps[0].CommentCount != 1 {
		t.Error("Expected comment count to be 1, got:", ps[0].CommentCount)
	}

	if ps[0].LikeCount != 1 {
		t.Error("Expected like count to be 1, got:", ps[0].LikeCount)
	}

	if !ps[0].Liked {
		t.Error("Expected post to be liked")
	}
}
