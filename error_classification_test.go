package toolsdk

import (
	"testing"

	"github.com/domainry/domainry-foundation/apperror"
)

func TestErrorClassificationContract(t *testing.T) {
	tests := []struct {
		class ErrorClass
		kind  apperror.ErrorKind
		code  string
	}{
		{ErrorClassBadRequest, apperror.KindBadRequest, "tools.validation_failed"},
		{ErrorClassForbidden, apperror.KindForbidden, "tools.permission_denied"},
		{ErrorClassNotFound, apperror.KindNotFound, "tools.resource_not_found"},
		{ErrorClassConflict, apperror.KindConflict, "tools.revision_conflict"},
		{ErrorClassUnavailable, apperror.KindUnavailable, "tools.dependency_unavailable"},
	}
	for _, test := range tests {
		err := &Error{Class: test.class, Code: test.code}
		if got := apperror.KindOf(err); got != test.kind {
			t.Fatalf("class=%q kind=%q", test.class, got)
		}
		if got := apperror.CodeOf(err); got != test.code {
			t.Fatalf("class=%q code=%q", test.class, got)
		}
	}
}
