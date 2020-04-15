package strutils

import "regexp"

var (
	PATTERN_URL   = regexp.MustCompile("^[a-zA-Z]+:\\/\\/[a-zA-Z0-9-_.]+[\\/]?(\\/[a-zA-Z0-9-_.#?%=&]+?)*$")
	PATTERN_EMAIL = regexp.MustCompile("^[a-zA-Z0-9\\-_\\.]+@[a-zA-Z0-9\\-_\\.]+\\.[a-zA-Z0-9]+$")
)

// IsURL returns whether the given string is a legal url
func IsURL(str string) bool {
	if len(str) < 3 {
		return false
	}
	return PATTERN_URL.Match([]byte(str))
}

// IsEmail returns whether the given string is a legal email address
func IsEmail(str string) bool {
	if len(str) < 3 {
		return false
	}
	return PATTERN_EMAIL.Match([]byte(str))
}
