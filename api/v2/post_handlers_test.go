package v2

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type postHandlerTestCase struct {
	method         string
	path           string
	body           io.Reader
	expectedStatus int
	expectedBody   string
}

func TestPostRoutes(t *testing.T) {
	var postTestCases = []postHandlerTestCase{
		{
			method:         "GET",
			path:           "/posts",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   "posts",
		},
		{
			method:         "GET",
			path:           "/posts/hailsatan",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			method:         "GET",
			path:           "/posts/" + testPostID,
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   "post",
		},
		{
			method:         "POST",
			path:           "/posts",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "POST",
			path:           "/posts",
			body:           strings.NewReader(`{"body":"hello","parent_type":"post","parent_id":"anything"}`),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "POST",
			path:           "/posts",
			body:           strings.NewReader(`{"parent_type":"post","parent_id":"something","body":"hello","action":"post"}`),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "POST",
			path:           "/posts",
			body:           strings.NewReader(`{"body":"hello","action":"post"}`),
			expectedStatus: http.StatusCreated,
			expectedBody:   "post",
		},
		{
			method:         "PATCH",
			path:           "/posts/" + testPostID,
			body:           strings.NewReader(`{"action":"like"}`),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "PATCH",
			path:           "/posts/" + unownedTestPostID,
			body:           strings.NewReader(`{"body":"hello"}`),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "PATCH",
			path:           "/posts/" + testPostID,
			body:           strings.NewReader(`{"body":"hello"}`),
			expectedStatus: http.StatusOK,
			expectedBody:   "post",
		},
		{
			method:         "DELETE",
			path:           "/posts/" + unownedTestPostID,
			body:           strings.NewReader(``),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			method:         "DELETE",
			path:           "/posts/" + testPostID,
			body:           strings.NewReader(``),
			expectedStatus: http.StatusOK,
			expectedBody:   "deleted",
		},
	}

	for _, tc := range postTestCases {
		w := performRequest(testRouter, tc.method, tc.path, tc.body)
		if w.Code != tc.expectedStatus {
			t.Error("Expected status to be:", tc.expectedStatus, "but got:", w.Code)
		}

		if tc.expectedBody != "" && !strings.Contains(w.Body.String(), tc.expectedBody) {
			t.Error("Expected body to contain:", tc.expectedBody, "but got:", w.Body.String())
		}
	}
}
