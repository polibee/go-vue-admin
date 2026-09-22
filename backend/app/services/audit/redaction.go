package auditservices

import (
	"encoding/json"
	"strings"
)

const auditPayloadLimitBytes = 32 * 1024

const redactedValue = "[REDACTED]"

var sensitiveAuditKeys = map[string]struct{}{
	"password": {}, "password_confirmation": {}, "token": {}, "access_token": {},
	"refresh_token": {}, "authorization": {}, "cookie": {}, "set-cookie": {},
	"secret": {}, "api_key": {}, "client_secret": {},
}

func RedactValue(value any) any {
	switch current := value.(type) {
	case map[string]any:
		redacted := make(map[string]any, len(current))
		for key, nested := range current {
			if _, sensitive := sensitiveAuditKeys[strings.ToLower(strings.TrimSpace(key))]; sensitive {
				redacted[key] = redactedValue
				continue
			}
			redacted[key] = RedactValue(nested)
		}
		return redacted
	case []any:
		redacted := make([]any, len(current))
		for index, nested := range current {
			redacted[index] = RedactValue(nested)
		}
		return redacted
	default:
		return value
	}
}

func MarshalBounded(value any) ([]byte, bool, error) {
	encoded, err := json.Marshal(RedactValue(value))
	if err != nil {
		return nil, false, err
	}
	if len(encoded) <= auditPayloadLimitBytes {
		return encoded, false, nil
	}
	bounded, err := json.Marshal(map[string]any{
		"truncated":      true,
		"original_bytes": len(encoded),
	})
	return bounded, true, err
}

func BoundedValue(value any) any {
	encoded, err := json.Marshal(RedactValue(value))
	if err != nil || len(encoded) <= auditPayloadLimitBytes {
		return RedactValue(value)
	}
	return map[string]any{"truncated": true, "original_bytes": len(encoded)}
}

func BoundedBytes(value []byte) ([]byte, bool) {
	if len(value) <= auditPayloadLimitBytes {
		return append([]byte(nil), value...), false
	}
	return append([]byte(nil), value[:auditPayloadLimitBytes]...), true
}
