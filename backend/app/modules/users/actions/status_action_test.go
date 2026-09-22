package actions

import (
	"errors"
	"testing"
)

func TestParseStatusParamsRejectsUnknownAndInvalidValues(t *testing.T) {
	if _, err := ParseStatusParams(map[string]any{"status": "disabled", "unexpected": true}); err != ErrUnknownParameter {
		t.Fatalf("expected unknown parameter error, got %v", err)
	}
	if _, err := ParseStatusParams(map[string]any{"status": "pending"}); err != ErrInvalidStatus {
		t.Fatalf("expected invalid status error, got %v", err)
	}
}

func TestParseStatusParamsRequiresStatus(t *testing.T) {
	if _, err := ParseStatusParams(map[string]any{}); err != ErrInvalidStatus {
		t.Fatalf("expected missing status error, got %v", err)
	}
}

func TestUserStatusErrorCodeDoesNotExposeInternalErrorTypes(t *testing.T) {
	if got := userStatusErrorCode(errors.New("database details")); got != "INTERNAL_ERROR" {
		t.Fatalf("unexpected internal error code: %s", got)
	}
}
