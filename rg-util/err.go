package rgutil

// Err represents a custom error type.
type Err struct {
	Message string
}

// NewErr creates a new Err instance with the given message.
func NewErr(message string) Err {
	return Err{Message: message}
}

// Error implements the error interface.
func (e Err) Error() string {
	return e.Message
}

// ErrNil represents a nil error for custom error handling.
var ErrNil = Err{Message: ""}
