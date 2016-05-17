package v2

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"partisan/models.v2/users"
)

func TestMain(m *testing.M) {
	defer testDB.Exec("DELETE FROM users;")

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
	m.Run()
}

type userHandlerTestCase struct {
	method         string
	path           string
	body           io.Reader
	expectedStatus int
	expectedBody   string
	loggedIn       bool
}

var userTestCases = []userHandlerTestCase{
	{
		method:         "POST",
		path:           "/users",
		body:           strings.NewReader(``),
		expectedStatus: http.StatusNotAcceptable,
		expectedBody:   "",
		loggedIn:       false,
	},
	{
		method:         "POST",
		path:           "/users",
		body:           strings.NewReader(`{"username":"user1","email":"user1@email.com","postal_code":"11233","password":"password","password_confirm":"password"}`),
		expectedStatus: http.StatusCreated,
		expectedBody:   "token",
		loggedIn:       false,
	},
	{
		method:         "POST",
		path:           "/users",
		body:           strings.NewReader(`{"username":"user2","email":"user2@email.com","postal_code":"11233","password":"foo","password_confirm":"bar"}`),
		expectedStatus: http.StatusNotAcceptable,
		expectedBody:   "",
		loggedIn:       false,
	},
	{
		method:         "GET",
		path:           "/users",
		body:           nil,
		expectedStatus: http.StatusOK,
		expectedBody:   "user",
		loggedIn:       true,
	},
	{
		method:         "PATCH",
		path:           "/users",
		body:           strings.NewReader(`{"gender":"unicorn"}`),
		expectedStatus: http.StatusOK,
		expectedBody:   "unicorn",
		loggedIn:       true,
	},
}

func TestUserRoutes(t *testing.T) {
	for _, tc := range userTestCases {
		w := performRequest(testRouter, tc.method, tc.path, tc.body)
		if w.Code != tc.expectedStatus {
			t.Error("Expected status to be:", tc.expectedStatus, "but got:", w.Code)
		}

		if tc.expectedBody != "" && !strings.Contains(w.Body.String(), tc.expectedBody) {
			t.Error("Expected body to contain:", tc.expectedBody, "but got:", w.Body.String())
		}
	}
}
