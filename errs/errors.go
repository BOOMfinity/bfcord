package errs

import (
	"errors"
	"fmt"
)

var (
	HTTPNotFound = HTTPErrorWithCode(NewHTTPError("resource does not exist"), 404)

	ItemNotFound = NewLibraryError("item does not exist in store")

	TooManyRetries = NewAPIError("too many retries")

	ResultMustBePointer = NewHTTPError("result must be a pointer")

	HTTPUnauthorized = HTTPErrorWithCode(NewHTTPError("invalid credentials"), 401)

	ModalValueNotFound = NewLibraryError("there is no modal component with given custom ID")
)

func IsNotFound(err error) bool {
	if errors.Is(err, HTTPNotFound) {
		return true
	}
	if errors.Is(err, ItemNotFound) {
		return true
	}
	return false
}

type LibraryError struct {
	Message string
}

func (e LibraryError) Error() string {
	return "bfcord: " + e.Message
}

type APIError struct {
	LibraryError
}

func (e APIError) Error() string {
	return "api: " + e.Message
}

type HTTPError struct {
	LibraryError
	Code int
}

func (e HTTPError) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("http(%v): %v", e.Code, e.Message)
	}
	return "http: " + e.Message
}

func NewAPIError(msg string) *APIError {
	return &APIError{LibraryError{msg}}
}

func NewLibraryError(msg string) *LibraryError {
	return &LibraryError{msg}
}

func NewHTTPError(msg string) *HTTPError {
	return &HTTPError{LibraryError: LibraryError{msg}}
}

func HTTPErrorWithCode(err *HTTPError, code int) *HTTPError {
	err.Code = code
	return err
}
