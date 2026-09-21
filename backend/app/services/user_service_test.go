package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateUserInput(t *testing.T) {
	tests := []struct {
		name            string
		userName        string
		email           string
		password        string
		requirePassword bool
		wantErr         error
	}{
		{name: "requires a password when creating", userName: "Admin", email: "admin@example.com", password: "", requirePassword: true, wantErr: ErrInvalidUser},
		{name: "accepts an empty password when updating", userName: "Admin", email: "admin@example.com", password: "", requirePassword: false},
		{name: "rejects a short password", userName: "Admin", email: "admin@example.com", password: "short", requirePassword: true, wantErr: ErrInvalidUser},
		{name: "rejects missing identity fields", userName: "", email: "admin@example.com", password: "password123", requirePassword: true, wantErr: ErrInvalidUser},
		{name: "accepts a valid payload", userName: " Admin ", email: " admin@example.com ", password: "password123", requirePassword: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserInput(tt.userName, tt.email, tt.password, tt.requirePassword)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateUserStatus(t *testing.T) {
	for _, status := range []string{"active", "disabled", "locked"} {
		require.NoError(t, validateUserStatus(status))
	}
	require.ErrorIs(t, validateUserStatus("pending"), ErrInvalidUser)
	require.Equal(t, userStatusActive, normalizeUserStatus(""))
}
