package v2

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type answerHandlerTestCase struct {
	method         string
	path           string
	body           io.Reader
	expectedStatus int
	expectedBody   string
}

func TestAnswerRoutes(t *testing.T) {
	var answerTestCases = []answerHandlerTestCase{
		{
			method:         "PATCH",
			path:           "/answers",
			body:           strings.NewReader(``),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "PATCH",
			path:           "/answers",
			body:           strings.NewReader(`{"map":[],"agree":true}`),
			expectedStatus: http.StatusNotAcceptable,
			expectedBody:   "",
		},
		{
			method:         "PATCH",
			path:           "/answers",
			body:           strings.NewReader(`{"map":[1,2,3],"agree":true}`),
			expectedStatus: http.StatusOK,
			expectedBody:   "updated",
		},
	}

	for _, tc := range answerTestCases {
		w := performRequest(testRouter, tc.method, tc.path, tc.body)
		if w.Code != tc.expectedStatus {
			t.Error("Expected status to be:", tc.expectedStatus, "but got:", w.Code)
		}

		if tc.expectedBody != "" && !strings.Contains(w.Body.String(), tc.expectedBody) {
			t.Error("Expected body to contain:", tc.expectedBody, "but got:", w.Body.String())
		}
	}
}
