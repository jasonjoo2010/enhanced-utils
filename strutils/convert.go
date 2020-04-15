package strutils

import (
	"strings"
)

// ToUnderline converts any upper case char into underscore and lower case one EXCEPT first one
func ToUnderscore(str string) string {
	b := strings.Builder{}
	cnt := 0
	for _, ch := range str {
		if ch >= 'A' && ch <= 'Z' {
			cnt++
		}
	}
	b.Grow(len(str) + cnt)
	for i, ch := range str {
		if ch >= 'A' && ch <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			ch += 'a' - 'A'
		} else if ch == ' ' {
			continue
		}
		b.WriteRune(ch)
	}
	return b.String()
}

// ToCamel converts any underscore into upper case, like 'apple_banana' =>  'appleBanana'
func ToCamel(str string) string {
	b := strings.Builder{}
	cnt := 0
	for _, ch := range str {
		if ch == '_' {
			cnt++
		}
	}
	b.Grow(len(str) - cnt)
	has_underscore := false
	for _, ch := range str {
		if ch == '_' {
			has_underscore = true
			continue
		} else if has_underscore && b.Len() > 0 && ch >= 'a' && ch <= 'z' {
			// to upper
			ch -= 'a' - 'A'
		}
		has_underscore = false
		if ch != ' ' {
			if ch >= 'A' && ch <= 'Z' && b.Len() == 0 {
				// first char lowcase
				ch += 'a' - 'A'
			}
			b.WriteRune(ch)
		}
	}
	return b.String()
}
