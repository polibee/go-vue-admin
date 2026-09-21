package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserPublicStatus(t *testing.T) {
	user := User{Status: "disabled"}

	public := user.Public()

	assert.Equal(t, "disabled", public["status"])
	assert.NotContains(t, public, "is_active")
}
