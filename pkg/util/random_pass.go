package util

import (
	"github.com/sethvargo/go-password/password"
)

func RandomPass(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
	return password.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)
}
