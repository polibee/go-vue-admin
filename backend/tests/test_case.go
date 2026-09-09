package tests

import (
	"os"

	"github.com/goravel/framework/testing"

	"goravel/bootstrap"
)

func init() {
	if os.Getenv("APP_ENV") == "" {
		_ = os.Setenv("APP_ENV", "local")
	}
	if os.Getenv("AUTH_BOOTSTRAP_PASSWORD") == "" {
		_ = os.Setenv("AUTH_BOOTSTRAP_PASSWORD", "test-only-password")
	}
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
