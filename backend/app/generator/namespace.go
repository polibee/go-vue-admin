package generator

import "fmt"

const defaultNamespace = "admin"

func normalizeNamespace(value string) (string, error) {
	if value == "" {
		return defaultNamespace, nil
	}
	if value != "admin" && value != "app" {
		return "", fmt.Errorf("unsupported API namespace %q", value)
	}
	return value, nil
}

func permissionName(namespace, resource, action string) string {
	return namespace + "." + resource + "." + action
}
