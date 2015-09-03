package validators

import (
	"fmt"
	"regexp"
	"time"
)

func NonBlank(s, tag string) error {
	if s == "" {
		return fmt.Errorf("%q cannot be blank", tag)
	}
	return nil
}

func Match(s, pattern, tag string) error {
	match, _ := regexp.MatchString(pattern, s)
	if !match {
		return fmt.Errorf("%q must match %q: %q", tag, pattern, s)
	}
	return nil
}

func Length(s string, min, max int, tag string) error {
	if len(s) < min {
		return fmt.Errorf("%q cannot be shorter than %d characters", tag, min)
	}

	if len(s) > max {
		return fmt.Errorf("%q cannot be longer than %d characters", tag, max)
	}

	return nil
}

func NonZero(i int, tag string) error {
	if i == 0 {
		return fmt.Errorf("%q cannot be zero", tag)
	}
	return nil
}

func NonEpoch(t time.Time, tag string) error {
	if t.IsZero() {
		return fmt.Errorf("%q cannot be cannot be January 1, 1970", tag)
	}
	return nil
}
