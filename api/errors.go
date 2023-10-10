package api

import (
	"errors"
	"fmt"
)

var (
	errUnknownHTTPCode = errors.New("unknown http code")
	errFieldValidation = errors.New("error while validating field")
)

// UnknownHTTPCodeError is returned when the http code can not be handled.
type UnknownHTTPCodeError struct {
	statusCode int
}

func (e *UnknownHTTPCodeError) Error() string {
	return fmt.Sprintf("%v : %v", errUnknownHTTPCode, e.statusCode)
}

func newUnknownHTTPCodeError(statusCode int) *UnknownHTTPCodeError {
	codeErr := &UnknownHTTPCodeError{
		statusCode: statusCode,
	}

	return codeErr
}

// FieldValidationError is returned when the expected format of a field is incorrect.
type FieldValidationError struct {
	wrappedErr error
}

func (e *FieldValidationError) Error() string {
	return fmt.Sprintf("%v : %v", errFieldValidation, e.wrappedErr)
}

func newFieldValidationError(validationErr error) *FieldValidationError {
	formatErr := &FieldValidationError{
		wrappedErr: validationErr,
	}

	return formatErr
}
