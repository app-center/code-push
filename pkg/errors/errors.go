package errors

import (
	"fmt"
	"strings"
)

// Error represents a WTF error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return e.string() }

func (e Error) String() string { return e.string() }

func (e Error) string() string {
	code := string(e)
	if strings.HasPrefix(code, "FA_") {
		return code
	} else {
		return fmt.Sprintf("FA_%s", code)
	}
}
