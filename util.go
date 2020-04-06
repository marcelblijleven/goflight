package goflight

import (
	"time"
)

func checkString(s string) (checked string, ok bool) {
	if s == "" {
		return s, false
	}

	return s, true
}

func checkTime(t time.Time) (checked time.Time, ok bool) {
	if t.IsZero() {
		return t, false
	}

	return t, true
}
