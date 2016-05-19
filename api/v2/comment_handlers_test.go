package v2

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type commentHandlerTestCase struct {
	method         string
	path           string
	body           io.Reader
	expectedStatus int
	expectedBody   string
}

func TestCommentRoutes(t *testing.T) {
	var commentTestCases = []commentHandlerTestCase{
		{
			method:         "GET",
			path:           "/posts/" + testCommentPostID + "/comments",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   testCommentID,
		},
		{
			method:         "GET",
			path:           "/posts/hailsatan/comments",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			method:         "POST",
			path:           "/posts/" + testCommentPostID + "/comments",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "POST",
			path:           "/posts/" + testCommentPostID + "/comments",
			body:           strings.NewReader(`{"body":"hello","action":"comment"}`),
			expectedStatus: http.StatusCreated,
			expectedBody:   "comment",
		},
	}

	for _, tc := range commentTestCases {
		w := performRequest(testRouter, tc.method, tc.path, tc.body)
		if w.Code != tc.expectedStatus {
			t.Error("Expected status to be:", tc.expectedStatus, "but got:", w.Code)
		}

		if tc.expectedBody != "" && !strings.Contains(w.Body.String(), tc.expectedBody) {
			t.Error("Expected body to contain:", tc.expectedBody, "but got:", w.Body.String())
		}
	}
}
