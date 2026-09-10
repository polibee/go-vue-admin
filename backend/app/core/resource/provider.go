package resource

import (
	"fmt"
	"strings"
)

type ResourceProviderMode string

const (
	ResourceProviderMemory   ResourceProviderMode = "memory"
	ResourceProviderDatabase ResourceProviderMode = "database"
	// ResourceProviderMySQL is kept as a compatibility alias for existing env files.
	ResourceProviderMySQL ResourceProviderMode = "mysql"
)

func ParseResourceProviderMode(value string) (ResourceProviderMode, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", string(ResourceProviderMemory):
		return ResourceProviderMemory, nil
	case string(ResourceProviderDatabase), string(ResourceProviderMySQL), "postgres", "pgsql":
		return ResourceProviderDatabase, nil
	default:
		return "", fmt.Errorf("unsupported resource provider %q", value)
	}
}
