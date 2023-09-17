package domain

import (
	"errors"
	"fmt"
)

const (
	// Error codes

	// NotFoundCode error
	NotFoundCode string = "EC100404"

	// DuplicateCode error
	DuplicateCode string = "EC100409"

	// ValidationCode error
	ValidationCode string = "EC100400"

	// DatabaseCode error
	DatabaseCode string = "EC100510"

	// NetworkCode error
	NetworkCode string = "EC100511"

	// TemplateCode error
	TemplateCode string = "EC100512"

	// MarshalCode error
	MarshalCode string = "EC100511"

	// InternalCode error
	InternalCode string = "EC100599"

	// AuthCode error
	AuthCode string = "EC100401"

	// MultiplePermissionsMatchedCode error
	MultiplePermissionsMatchedCode string = "EC100451"

	// ConflictingPermissionsCode error
	ConflictingPermissionsCode string = "EC100452"
)

// ValidationError error
type ValidationError struct {
	Message string
}

// NewValidationError constructor
func NewValidationError(msg string) *ValidationError {
	return &ValidationError{
		Message: msg + " [" + ValidationCode + "]",
	}
}

// Error getter
func (e *ValidationError) Error() string {
	return e.Message
}

// String getter
func (e *ValidationError) String() string {
	return fmt.Sprintf("ValidationError: %s", e.Message)
}

// DatabaseError error
type DatabaseError struct {
	Message string
}

// NewDatabaseError constructor
func NewDatabaseError(msg string) *DatabaseError {
	return &DatabaseError{
		Message: msg + " [" + DatabaseCode + "]",
	}
}

func (e *DatabaseError) Error() string {
	return e.Message
}

// String getter
func (e *DatabaseError) String() string {
	return fmt.Sprintf("DatabaseError: %s", e.Message)
}

// MarshalError error
type MarshalError struct {
	Message string
}

// NewMarshalError constructor
func NewMarshalError(msg string) *MarshalError {
	return &MarshalError{
		Message: msg + " [" + MarshalCode + "]",
	}
}

func (e *MarshalError) Error() string {
	return e.Message
}

// String getter
func (e *MarshalError) String() string {
	return fmt.Sprintf("MarshalError: %s", e.Message)
}

// DuplicateError error
type DuplicateError struct {
	Message string
}

// NewDuplicateError constructor
func NewDuplicateError(msg string) *DuplicateError {
	return &DuplicateError{
		Message: msg + " [" + DuplicateCode + "]",
	}
}

func (e *DuplicateError) Error() string {
	return e.Message
}

// String getter
func (e *DuplicateError) String() string {
	return fmt.Sprintf("DuplicateError: %s", e.Message)
}

// NotFoundError error
type NotFoundError struct {
	Message string
}

// NewNotFoundError constructor
func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		Message: msg + " [" + NotFoundCode + "]",
	}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// String getter
func (e *NotFoundError) String() string {
	return fmt.Sprintf("NotFoundError: %s", e.Message)
}

// AuthError error
type AuthError struct {
	Message string
}

// NewAuthError constructor
func NewAuthError(msg string) *AuthError {
	return &AuthError{
		Message: msg + " [" + AuthCode + "]",
	}
}

// NewAuthErrorWithCode constructor
func NewAuthErrorWithCode(msg string, code string) *AuthError {
	return &AuthError{
		Message: msg + " [" + code + "]",
	}
}

func (e *AuthError) Error() string {
	return e.Message
}

// String getter
func (e *AuthError) String() string {
	return fmt.Sprintf("AuthError: %s", e.Message)
}

// InternalError error
type InternalError struct {
	Message string
}

// NewInternalError constructor
func NewInternalError(msg string, code string) *InternalError {
	return &InternalError{
		Message: msg + " [" + code + "]",
	}
}

func (e *InternalError) Error() string {
	return e.Message
}

// String getter
func (e *InternalError) String() string {
	return fmt.Sprintf("InternalError: %s", e.Message)
}

// ErrorToHTTPStatus helper
func ErrorToHTTPStatus(err error) int {
	var validationErr *ValidationError
	var marshalError *MarshalError
	var notFoundErr *NotFoundError
	var duplicateErr *DuplicateError
	var authFoundErr *AuthError
	var databaseErr *DatabaseError
	var internalErr *InternalError
	if errors.As(err, &validationErr) {
		return 400
	} else if errors.As(err, &marshalError) {
		return 400
	} else if errors.As(err, &notFoundErr) {
		return 404
	} else if errors.As(err, &duplicateErr) {
		return 409
	} else if errors.As(err, &authFoundErr) {
		return 401
	} else if errors.As(err, &databaseErr) {
		return 500
	} else if errors.As(err, &internalErr) {
		return 500
	} else if err != nil {
		return 500
	}
	return 200
}
