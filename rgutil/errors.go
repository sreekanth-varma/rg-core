package rgutil

import (
	"fmt"
	"os"
	"runtime/debug"
)

// // Err represents a custom error type.
// type Err struct {
// 	Message string
// }

// // NewErr creates a new Err instance with the given message.
// func NewErr(message string) Err {
// 	return Err{Message: message}
// }

// // Error implements the error interface.
// func (e Err) Error() string {
// 	return e.Message
// }

// // ErrNil represents a nil error for custom error handling.
// var ErrNil = Err{Message: ""}

type Err uint16

const (
	// 200: success
	ErrNil Err = 200

	// 201: success, but no data
	ErrNoData Err = 201

	// 400: bad input
	ErrBadInput Err = 400

	// 401: not authenticated
	ErrNotAuthenticated Err = 401

	// 403: not authorised
	ErrNotAuthorised Err = 403

	// 422: business validation failed
	ErrNotValid Err = 422

	// 500: server error. cannot retry
	ErrProcessingFailed Err = 500

	// 503: server error. can retry
	ErrUnavailable Err = 503
)

func PanicHandler() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "Recovered. Error: %v\n", err)
		fmt.Println(string(debug.Stack()))
		return
	}
}
