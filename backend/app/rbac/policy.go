package rbac

import (
	"errors"
	"strings"
)

var ErrInvalidRoleIDs = errors.New("role ids must be positive and unique")
var ErrInvalidRoleInput = errors.New("role name and display name are required")

func HasAnyPermission(permissions map[string]struct{}, required ...string) bool {
	for _, name := range required {
		if _, ok := permissions[name]; ok {
			return true
		}
	}
	return false
}

func ValidateRoleIDs(roleIDs []int64) error {
	seen := make(map[int64]struct{}, len(roleIDs))
	for _, roleID := range roleIDs {
		if roleID <= 0 {
			return ErrInvalidRoleIDs
		}
		if _, exists := seen[roleID]; exists {
			return ErrInvalidRoleIDs
		}
		seen[roleID] = struct{}{}
	}
	return nil
}

func ValidateRoleInput(name, displayName string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(displayName) == "" {
		return ErrInvalidRoleInput
	}
	return nil
}
