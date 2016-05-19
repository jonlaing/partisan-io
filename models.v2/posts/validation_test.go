package posts

import (
	"database/sql"
	"testing"

	"github.com/nu7hatch/gouuid"

	models "partisan/models.v2"
)

type validationTest struct {
	post           Post
	errorLength    int
	errorsExpected models.ValidationErrors
}

var testuuid, _ = uuid.NewV4()

var tests = []validationTest{
	{
		post:           Post{UserID: "not a uuid", Action: APost},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"user_id": models.ErrUUIDFormat},
	},
	{
		post:           Post{UserID: testuuid.String(), ParentID: sql.NullString{"not a uuid", true}, ParentType: PTPost, Action: AComment},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"parent_id": models.ErrUUIDFormat},
	},
	{
		post:           Post{UserID: testuuid.String(), ParentID: sql.NullString{testuuid.String(), true}, Action: APost, ParentType: PTPost},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"parent_type": ErrParentType},
	},
	{
		post:           Post{UserID: testuuid.String(), ParentID: sql.NullString{testuuid.String(), true}, Action: Action("something else")},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"action": ErrAction},
	},
	{
		post:           Post{UserID: testuuid.String(), ParentID: sql.NullString{testuuid.String(), true}, Action: APost, ParentType: ParentType("something else")},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"parent_type": ErrParentType},
	},
	{
		post:           Post{UserID: testuuid.String(), ParentID: sql.NullString{testuuid.String(), true}, Action: ALike, ParentType: PTGroup},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"parent_type": ErrLikeParent},
	},
	{
		post:           Post{UserID: testuuid.String(), ParentID: sql.NullString{testuuid.String(), true}, Action: ALike, ParentType: PTPost, Body: "hello"},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"body": ErrLikeBody},
	},
}

func TestValidation(t *testing.T) {
	for _, test := range tests {
		errs := test.post.Validate()
		if len(errs) != test.errorLength {
			t.Error("Expected", test.errorLength, "errors, got", len(errs), ":", test.errorsExpected, "--", errs)
		}

		for k := range test.errorsExpected {
			if _, ok := errs[k]; !ok {
				t.Error("Expected", k, "to be in list of errors. Got:", errs)
			}

			if errs[k] != test.errorsExpected[k] {
				t.Error("Expected error on", k, "to be:", test.errorsExpected[k], "Got:", errs[k])
			}

		}
	}
}
