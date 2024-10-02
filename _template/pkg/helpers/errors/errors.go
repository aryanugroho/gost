package errors

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/flip-id/go-core/helpers/method"
)

// ErrorType is the type of an error
type ErrorType uint

// NoType error
const (
	NoType              ErrorType = 500
	BadRequest          ErrorType = 400
	Unauthorized        ErrorType = 401
	Forbidden           ErrorType = 403
	NotFound            ErrorType = 404
	Validation          ErrorType = 422
	UnprocessableEntity ErrorType = 422
)

type ErrorCode uint

const (
	NoCode                    ErrorCode = 0
	AuthorizationFailedCode   ErrorCode = 401
	APINotFoundFailedCode     ErrorCode = 404
	InternalServerErrorCode   ErrorCode = 500
	MaintenanceStatusCode     ErrorCode = 503
	InvalidInputParameterCode ErrorCode = 4001
)

var (
	ErrorInvalidParam = BadRequest.New("Invalid parameter", InvalidInputParameterCode)
)

type CustomError struct {
	errorType     ErrorType
	errorCode     ErrorCode
	originalError error
	contexts      []errorContext
	stackTrace    []string
}

type errorContext struct {
	Field   string
	Message string
}

// New creates a new CustomError
func (errorType ErrorType) New(msg string, codes ...ErrorCode) error {
	var code ErrorCode
	if len(codes) > 0 {
		code = codes[0]
	}
	customErr := CustomError{errorType: errorType, errorCode: code, originalError: errors.New(msg), stackTrace: []string{msg}}
	return customErr
}

// Newf creates a new CustomError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	customErr := CustomError{errorType: errorType, originalError: fmt.Errorf(msg, args...), stackTrace: []string{msg}}
	return customErr
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string, codes ...ErrorCode) error {
	var code ErrorCode
	if len(codes) > 0 {
		code = codes[0]
	}
	return CustomError{errorType: errorType, errorCode: code, originalError: errors.Wrapf(err, msg)}
}

// Wrapf creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return CustomError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

// Error returns the mssage of a CustomError
func (error CustomError) Error() string {
	return error.originalError.Error()
}

// New creates a no type error and report to sentry
func New(msg string) error {
	err := CustomError{errorType: NoType, originalError: errors.New(msg)}
	return err
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	err := CustomError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
	return err
}

func Msg(err error, msg string) error {
	fileName, line, fnName := method.TraceCaller(3)
	errorMsg := fmt.Sprintf("%s:%d@%s()", fileName, line, fnName)
	if customErr, ok := err.(CustomError); ok {
		return CustomError{
			errorType:     customErr.errorType,
			originalError: errors.New(msg),
			contexts:      customErr.contexts,
			stackTrace:    append([]string{errorMsg}, customErr.stackTrace...),
		}
	}

	return CustomError{errorType: NoType, originalError: errors.New(msg), stackTrace: []string{errorMsg}}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(CustomError); ok {
		return CustomError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			contexts:      customErr.contexts,
			stackTrace:    customErr.stackTrace,
		}
	}

	return CustomError{errorType: NoType, originalError: wrappedError}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// AddStackTrace an error with format string
func AddStackTrace(err error, msg string) error {
	if customErr, ok := err.(CustomError); ok {
		stackTrace := append([]string{msg}, customErr.stackTrace...)
		return CustomError{errorType: customErr.errorType, originalError: customErr.originalError, contexts: customErr.contexts, stackTrace: stackTrace}
	}

	stackTrace := []string{msg}
	return CustomError{errorType: NoType, originalError: err, stackTrace: stackTrace}
}

// GetStackTrace returns the error stack trace
func GetStackTrace(err error) []string {
	if customErr, ok := err.(CustomError); ok {
		return customErr.stackTrace
	}
	return []string{}
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(CustomError); ok {
		return customErr.errorType
	}

	return NoType
}

func GetErrorCode(err error) int {
	if customErr, ok := err.(CustomError); ok {
		if customErr.errorCode == NoCode {
			return int(customErr.errorType)
		}
		return int(customErr.errorCode)
	}

	return int(InternalServerErrorCode)
}

// Is Check if error is the specified error type
func Is(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}

	errType := NoType
	if customErr, ok := err.(CustomError); ok {
		errType = customErr.errorType
	}

	return errType == errorType
}

// IsNotFound Check if error is NotFound error type
func IsNotFound(err error) bool {
	return Is(err, NotFound)
}
