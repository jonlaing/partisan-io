package validators

import (
	"testing"
	"time"
)

func TestNonBlank(t *testing.T) {
	err := NonBlank("blah", "test_tag")
	if err != nil {
		t.Fail()
	}

	err = NonBlank("", "test_tag")
	if err == nil {
		t.Fail()
	}

	if err.Field != "test_tag" {
		t.Fail()
	}
}

func TestMatch(t *testing.T) {
	err := Match("blah", "bl.+", "test_tag")
	if err != nil {
		t.Fail()
	}

	err = Match("blah", "o+", "test_tag")
	if err == nil {
		t.Fail()
	}

	if err.Field != "test_tag" {
		t.Fail()
	}
}

func TestLength(t *testing.T) {
	err := Length("123", 0, 3, "test_tag")
	if err != nil {
		t.Fail()
	}

	err = Length("123", 4, 5, "test_tag")
	if err == nil {
		t.Fail()
	}

	if err.Field != "test_tag" {
		t.Fail()
	}
}

func TestNonZero(t *testing.T) {
	err := NonZero(1, "test_tag")
	if err != nil {
		t.Fail()
	}

	err = NonZero(0, "test_tag")
	if err == nil {
		t.Fail()
	}

	if err.Field != "test_tag" {
		t.Fail()
	}
}

func TestNonEpoch(t *testing.T) {
	err := NonEpoch(time.Now(), "test_tag")
	if err != nil {
		t.Fail()
	}

	err = NonEpoch(time.Time{}, "test_tag")
	if err == nil {
		t.Fail()
	}

	if err.Field != "test_tag" {
		t.Fail()
	}
}
