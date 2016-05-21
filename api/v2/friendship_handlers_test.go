package v2

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type friendshipHandlerTestCase struct {
	method         string
	path           string
	body           io.Reader
	expectedStatus int
	expectedBody   string
}

func TestFriendshipRoutes(t *testing.T) {
	var friendshipTestCases = []friendshipHandlerTestCase{
		{
			method:         "GET",
			path:           "/friendships",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   testFriendID,
		},
		{
			method:         "GET",
			path:           "/friendships/hailsatan",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			method:         "GET",
			path:           "/friendships/" + testFriendID,
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   testFriendID,
		},
		{
			method:         "POST",
			path:           "/friendships",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "POST",
			path:           "/friendships",
			body:           strings.NewReader(fmt.Sprintf(`{"friend_id":"%s"}`, testUnfriendedID)),
			expectedStatus: http.StatusCreated,
			expectedBody:   testUnfriendedID,
		},
		{
			method:         "PATCH",
			path:           "/friendships/" + testUnfriendedID,
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "PATCH",
			path:           "/friendships/" + testUnconfirmedID,
			body:           strings.NewReader(`{"confirmed":true}`),
			expectedStatus: http.StatusOK,
			expectedBody:   testUnconfirmedID,
		},
		{
			method:         "DELETE",
			path:           "/friendships/" + testFriendID,
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   "deleted",
		},
	}

	for _, tc := range friendshipTestCases {
		w := performRequest(testRouter, tc.method, tc.path, tc.body)
		if w.Code != tc.expectedStatus {
			t.Error(tc.method, tc.path, "Expected status to be:", tc.expectedStatus, "but got:", w.Code)
		}

		if tc.expectedBody != "" && !strings.Contains(w.Body.String(), tc.expectedBody) {
			t.Error(tc.method, tc.path, "Expected body to contain:", tc.expectedBody, "but got:", w.Body.String())
		}
	}
}
