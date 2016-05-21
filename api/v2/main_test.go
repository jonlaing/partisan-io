package v2

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"partisan/auth"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"
	"partisan/models.v2/friendships"
	"partisan/models.v2/posts"
	"partisan/models.v2/users"
)

var testDB *gorm.DB
var testRouter *gin.Engine
var userCount = 0
var testPostID, unownedTestPostID string
var testLikePostID string
var testCommentPostID, testCommentID string
var testFriendID, testUnconfirmedID, testUnfriendedID string

func init() {
	var err error
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PW"))
	testDB, err = gorm.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	if err := testDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		panic(err)
	}

	testDB.AutoMigrate(users.User{}, posts.Post{}, friendships.Friendship{})

	testRouter = gin.Default()
	testRouter.Use(addTestDB())
}

func TestMain(m *testing.M) {
	defer testDB.Exec("DELETE FROM users;")
	defer testDB.Exec("DELETE FROM posts;")
	defer testDB.Exec("DELETE FROM friendships;")
	initUserTests()
	testPostID, unownedTestPostID = initPostTests()
	testLikePostID = initLikeTests()
	testCommentPostID, testCommentID = initCommentTests()
	testFriendID, testUnconfirmedID, testUnfriendedID = initFriendshipTests()
	m.Run()
}

func addTestDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := testDB.DB().Ping(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Set("db", testDB)
	}
}

func createTestUser() users.User {
	userCount++
	name := fmt.Sprintf("user_%d", userCount)
	uBinding := users.CreatorBinding{
		Username:        name,
		Email:           name + "@email.com",
		PostalCode:      "11233",
		Password:        "password",
		PasswordConfirm: "password",
	}

	testUser, errs := users.New(uBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	testUser.GenAPIKey()

	if err := testDB.Save(&testUser).Error; err != nil {
		panic(err)
	}

	return testUser
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func login(user *users.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		tok, _ := auth.Login(user, c)
		c.Request.Header.Set("X-Auth-Token", tok)
		auth.Auth()(c)
	}
}

func initUserTests() {
	uBinding := users.CreatorBinding{
		Username:        "user",
		Email:           "user@email.com",
		PostalCode:      "11233",
		Password:        "password",
		PasswordConfirm: "password",
	}

	testUser, errs := users.New(uBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	testUser.GenAPIKey()

	if err := testDB.Save(&testUser).Error; err != nil {
		panic(err)
	}

	testRouter.POST("/users", UserCreate)
	testRouter.GET("/users", login(&testUser), UserShow) // Show Current User
	testRouter.PATCH("/users", login(&testUser), UserUpdate)
	// TODO: test this one...
	// testRouter.POST("/users/avatar_upload", login(&testUser), UserAvatarUpload)
}

func initPostTests() (string, string) {
	user := createTestUser()
	otherUserID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	pBinding := posts.CreatorBinding{
		Body:   "my post",
		Action: posts.APost,
	}

	testPost, errs := posts.New(user.ID, pBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	if err := testDB.Save(&testPost).Error; err != nil {
		panic(err)
	}

	puBinding := posts.CreatorBinding{
		Body:   "unowned post",
		Action: posts.APost,
	}

	unownedPost, errs := posts.New(otherUserID.String(), puBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	if err := testDB.Save(&unownedPost).Error; err != nil {
		panic(err)
	}

	testRouter.GET("/posts", login(&user), PostIndex)
	testRouter.POST("/posts", login(&user), PostCreate)
	testRouter.GET("/posts/:record_id", login(&user), PostShow)
	testRouter.PATCH("/posts/:record_id", login(&user), PostUpdate)
	testRouter.DELETE("/posts/:record_id", login(&user), PostDestroy)

	return testPost.ID, unownedPost.ID
}

func initLikeTests() string {
	user := createTestUser()

	pBinding := posts.CreatorBinding{
		Body:   "my post",
		Action: posts.APost,
	}

	testLikePost, errs := posts.New(user.ID, pBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	if err := testDB.Save(&testLikePost).Error; err != nil {
		panic(err)
	}

	testRouter.POST("/posts/:record_id/like", login(&user), LikeCreate)

	return testLikePost.ID
}

func initCommentTests() (string, string) {
	user := createTestUser()

	pBinding := posts.CreatorBinding{
		Body:   "my post",
		Action: posts.APost,
	}

	testPost, errs := posts.New(user.ID, pBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	if err := testDB.Save(&testPost).Error; err != nil {
		panic(err)
	}

	cBinding := posts.CreatorBinding{
		ParentID:   testPost.ID,
		ParentType: string(testPost.PostParentType()),
		Body:       "a comment",
		Action:     posts.AComment,
	}

	comment, errs := posts.New(user.ID, cBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	if err := testDB.Save(&comment).Error; err != nil {
		panic(err)
	}

	testRouter.GET("/posts/:record_id/comments", login(&user), CommentIndex)
	testRouter.POST("/posts/:record_id/comments", login(&user), CommentCreate)

	return testPost.ID, comment.ID
}

func initFriendshipTests() (string, string, string) {
	user := createTestUser()  // Me
	user2 := createTestUser() // Friended
	user3 := createTestUser() // Unconfirmed
	user4 := createTestUser() // Unfriended

	fBinding := friendships.CreatorBinding{
		FriendID: user2.ID,
	}

	testFriendship1, errs := friendships.New(user.ID, fBinding)
	if len(errs) > 0 {
		panic(errs)
	}

	testFriendship1.Confirmed = true

	if err := testDB.Save(&testFriendship1).Error; err != nil {
		panic(err)
	}

	f2Binding := friendships.CreatorBinding{
		FriendID: user.ID,
	}

	testFriendship2, errs := friendships.New(user3.ID, f2Binding)
	if len(errs) > 0 {
		panic(errs)
	}

	if err := testDB.Save(&testFriendship2).Error; err != nil {
		panic(err)
	}

	testRouter.GET("/friendships", login(&user), FriendshipIndex)
	testRouter.GET("/friendships/:user_id", login(&user), FriendshipShow)
	testRouter.POST("/friendships", login(&user), FriendshipCreate)
	testRouter.PATCH("/friendships/:user_id", login(&user), FriendshipUpdate)
	testRouter.DELETE("/friendships/:user_id", login(&user), FriendshipDestroy)

	return user2.ID, user3.ID, user4.ID
}
