package v2

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type userHandlerTestCase struct {
	method         string
	path           string
	body           io.Reader
	expectedStatus int
	expectedBody   string
	loggedIn       bool
}

func TestUserRoutes(t *testing.T) {
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
			method:         "GET",
			path:           "/users/" + testUser.Username,
			body:           nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "match",
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
