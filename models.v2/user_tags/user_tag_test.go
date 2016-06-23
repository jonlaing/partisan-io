package usertags

import "testing"

type tagger struct {
	body string
}

func (t tagger) GetID() string {
	return "id"
}

func (t tagger) GetUserID() string {
	return "userid"
}

func (t tagger) GetType() string {
	return "tagger"
}

func (t tagger) GetContent() string {
	return t.body
}

type testCase struct {
	body           string
	expectedNumber int
	expectedTags   []string
}

var testCases = []testCase{
	{
		body:           "this doesn't have a tag in it",
		expectedNumber: 0,
		expectedTags:   []string{},
	},
	{
		body:           "this has a @user tag in it",
		expectedNumber: 1,
		expectedTags:   []string{"user"},
	},
	{
		body:           "this has two @user @blah tags in it",
		expectedNumber: 2,
		expectedTags:   []string{"user", "blah"},
	},
	{
		body:           "this has an email but no tag user@email.com",
		expectedNumber: 0,
		expectedTags:   []string{},
	},
	{
		body:           "this has a @user tag and an email blah@email.com",
		expectedNumber: 1,
		expectedTags:   []string{"user"},
	},
}

func TestExtractTags(t *testing.T) {
	for _, tc := range testCases {
		r := tagger{body: tc.body}
		tags := extractTags(r.GetContent())
		if len(tags) != tc.expectedNumber {
			t.Error("Expected", tc.expectedNumber, "but got:", len(tags))
		}

		for _, eTag := range tc.expectedTags {
			found := false
			for _, tag := range tags {
				if eTag == tag {
					found = true
				}
			}

			if !found {
				t.Error("Expected to find", eTag, "in tags:", tags)
			}
		}
	}
}

func TestExtract(t *testing.T) {
}
