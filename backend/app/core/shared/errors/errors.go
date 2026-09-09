package errors

import "fmt"

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ErrorEnvelope struct {
	Error APIError `json:"error"`
}

func New(code, message string, details any) ErrorEnvelope {
	return ErrorEnvelope{Error: APIError{Code: code, Message: message, Details: details}}
}

func InvalidQuery(err error) ErrorEnvelope {
	return New("INVALID_QUERY", "请求查询参数无效", map[string]string{"reason": fmt.Sprint(err)})
}
