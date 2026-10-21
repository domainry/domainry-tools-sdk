package toolsdk

import (
	"strings"

	"github.com/domainry/domainry-foundation/apperror"
)

type ErrorClass = apperror.ErrorKind

const (
	ErrorClassBadRequest  = apperror.KindBadRequest
	ErrorClassForbidden   = apperror.KindForbidden
	ErrorClassNotFound    = apperror.KindNotFound
	ErrorClassConflict    = apperror.KindConflict
	ErrorClassRateLimited = apperror.KindRateLimited
	ErrorClassUnavailable = apperror.KindUnavailable
	ErrorClassInternal    = apperror.KindInternal
)

type Error struct {
	Class     ErrorClass
	Code      string
	Message   string
	Retryable bool
	Cause     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if strings.TrimSpace(e.Message) != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return strings.TrimSpace(e.Code)
}
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *Error) ErrorCode() string {
	if e == nil {
		return ""
	}
	return strings.TrimSpace(e.Code)
}

func (e *Error) ApplicationErrorKind() string {
	if e == nil {
		return ""
	}
	if e.Class == "" && e.Retryable {
		return string(apperror.KindUnavailable)
	}
	return string(apperror.NormalizeKind(string(e.Class)))
}

func (*Error) ErrorParams() map[string]string { return nil }
