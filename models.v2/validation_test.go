package modelsv2

import (
	"errors"
	"testing"
)

type validationTest struct {
	in  ValidationErrors
	out string
}

var vtests = []validationTest{
	validationTest{
		in:  ValidationErrors{"field": errors.New("Error 1")},
		out: "{\"field\":\"Error 1\"}",
	},
	validationTest{
		in:  ValidationErrors{"field": errors.New("Error 1"), "field2": errors.New("Error 2")},
		out: "{\"field\":\"Error 1\",\"field2\":\"Error 2\"}",
	},
}

func TestValidationErrors(t *testing.T) {
	for _, test := range vtests {
		if test.in.Error() != test.out {
			t.Error("Expected:", test.out, ", but got:", test.in.Error())
		}
	}
}
