package errors

// Error represents a WTF error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }

func (e Error) String() string { return string(e) }
