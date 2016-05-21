package friendships

import (
	"testing"

	"github.com/nu7hatch/gouuid"
	models "partisan/models.v2"
)

var testID, testID2 string

func init() {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	testID = id.String()

	id, err = uuid.NewV4()
	if err != nil {
		panic(err)
	}

	testID2 = id.String()
}

type validationTest struct {
	friendship     Friendship
	errorLength    int
	errorsExpected models.ValidationErrors
}

var tests = []validationTest{
	validationTest{
		friendship: Friendship{
			UserID:   "blah blah",
			FriendID: "hail satan",
		},
		errorLength:    2,
		errorsExpected: models.ValidationErrors{"user_id": models.ErrUUIDFormat, "friend_id": models.ErrUUIDFormat},
	},
	validationTest{
		friendship: Friendship{
			UserID:   testID,
			FriendID: testID,
		},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"friend_id": ErrFriendSelf},
	},
	validationTest{
		friendship: Friendship{
			UserID:   testID,
			FriendID: testID2,
		},
		errorLength:    0,
		errorsExpected: models.ValidationErrors{},
	},
}

func TestValidation(t *testing.T) {
	for _, test := range tests {
		errs := test.friendship.Validate()
		if len(errs) != test.errorLength {
			t.Error("Expected", test.errorLength, "errors, got", len(errs), ":", test)
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
