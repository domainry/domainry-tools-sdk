package toolsdk

import "strings"

type Error struct {
	Class, Code, Message string
	Retryable            bool
	Cause                error
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
	switch strings.ToLower(strings.TrimSpace(e.Class)) {
	case "bad_request", "invalid":
		return "bad_request"
	case "forbidden", "unauthenticated", "unauthorized":
		return "forbidden"
	case "not_found":
		return "not_found"
	case "conflict":
		return "conflict"
	case "rate_limited":
		return "rate_limited"
	case "unavailable":
		return "unavailable"
	default:
		if e.Retryable {
			return "unavailable"
		}
		return "internal"
	}
}

func (*Error) ErrorParams() map[string]string { return nil }
