package v2

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type likeHandlerTestCase struct {
	method         string
	path           string
	body           io.Reader
	expectedStatus int
	expectedBody   string
}

func TestLikeRoutes(t *testing.T) {
	var likeTestCases = []likeHandlerTestCase{
		{
			method:         "POST",
			path:           "/posts/" + testLikePostID + "/like",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusCreated,
			expectedBody:   "\"like_count\":1",
		},
		{
			method:         "POST",
			path:           "/posts/" + testLikePostID + "/like",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   "\"like_count\":0",
		},
	}

	for _, tc := range likeTestCases {
		w := performRequest(testRouter, tc.method, tc.path, tc.body)
		if w.Code != tc.expectedStatus {
			t.Error("Expected status to be:", tc.expectedStatus, "but got:", w.Code)
		}

		if tc.expectedBody != "" && !strings.Contains(w.Body.String(), tc.expectedBody) {
			t.Error("Expected body to contain:", tc.expectedBody, "but got:", w.Body.String())
		}
	}
}
