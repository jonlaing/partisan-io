package users

import (
	"testing"

	models "partisan/models.v2"
)

type validationTest struct {
	user           User
	errorLength    int
	errorsExpected models.ValidationErrors
}

var tests = []validationTest{
	validationTest{
		user: User{
			Username: "$% !",
			Email:    "hello@email.com",
		},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"username": ErrUsernameValidation},
	},
	validationTest{
		user: User{
			Username: "00000000000000000",
			Email:    "hello@email.com",
		},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"username": ErrUsernameValidation},
	},
	validationTest{
		user: User{
			Username: "0",
			Email:    "hello@email.com",
		},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"username": ErrUsernameValidation},
	},
	validationTest{
		user: User{
			Username: "username",
			Email:    "asoidhfoih",
		},
		errorLength:    1,
		errorsExpected: models.ValidationErrors{"email": ErrEmailValidation},
	},
	validationTest{
		user: User{
			Username: "!@ %",
			Email:    "asoidhfoih",
		},
		errorLength: 2,
		errorsExpected: models.ValidationErrors{
			"username": ErrUsernameValidation,
			"email":    ErrEmailValidation,
		},
	},
	validationTest{
		user: User{
			Username: "username",
			Email:    "hello@email.com",
		},
		errorLength:    0,
		errorsExpected: models.ValidationErrors{},
	},
}

func TestValidation(t *testing.T) {
	for _, test := range tests {
		errs := test.user.Validate()
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
