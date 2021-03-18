package strutils

import (
	"math"
	"math/rand"
)

const (
	CHARS_LOWCASE = "abcdefghijklmnopqrstuvwxyz" +
		"0123456789"
	CHARS_NORMAL = CHARS_LOWCASE +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CHARS_HASH = "abcdef" +
		"0123456789"
	CHARS_PRINTABLE = CHARS_NORMAL +
		"!\"#$%&'()*+,-./:;<=>?@" +
		"[\\]^_`{|}~"
)

// RandNumbers returns random number characters string in specific length
func RandNumbers(length int) string {
	num := RandUint64(length)
	arr := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		arr[i] = '0' + byte(num%10)
		num /= 10
	}
	return string(arr)
}

// RandUint64 returns a uint64 random number with max specified width
func RandUint64(max_width int) uint64 {
	return rand.Uint64() % uint64(math.Pow10(max_width))
}

func randString(length int, available string) string {
	mod := uint64(len(available))
	r := rand.Uint64()
	arr := make([]byte, length)
	for i := 0; i < length; i++ {
		arr[i] = available[r%mod]

		if r < mod {
			r = rand.Uint64()
		} else {
			r /= mod
		}
	}
	return string(arr)
}

// RandString returns random string including a-z, A-Z, 0-9
func RandString(length int) string {
	return randString(length, CHARS_NORMAL)
}

// RandLowCased returns random string including a-z, 0-9
func RandLowCased(length int) string {
	return randString(length, CHARS_LOWCASE)
}

// RandHash returns random string including a-f, 0-9
func RandHash(length int) string {
	return randString(length, CHARS_HASH)
}

// RandPrintable returns random string including special characters
func RandPrintable(length int) string {
	return randString(length, CHARS_PRINTABLE)
}
