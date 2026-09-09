package resource

import (
	"fmt"
	"strings"
)

type ResourceProviderMode string

const (
	ResourceProviderMemory ResourceProviderMode = "memory"
	ResourceProviderMySQL  ResourceProviderMode = "mysql"
)

func ParseResourceProviderMode(value string) (ResourceProviderMode, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", string(ResourceProviderMemory):
		return ResourceProviderMemory, nil
	case string(ResourceProviderMySQL):
		return ResourceProviderMySQL, nil
	default:
		return "", fmt.Errorf("unsupported resource provider %q", value)
	}
}
