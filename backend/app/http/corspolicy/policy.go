package corspolicy

import "strings"

var allowedOrigins = map[string]struct{}{
	"http://127.0.0.1:4173": {},
	"http://127.0.0.1:4174": {},
	"http://127.0.0.1:5173": {},
	"http://127.0.0.1:5174": {},
	"http://127.0.0.1:5175": {},
	"http://127.0.0.1:5176": {},
	"http://127.0.0.1:5180": {},
}

func AllowedOrigin(origin string) bool {
	_, ok := allowedOrigins[strings.TrimSpace(origin)]
	return ok
}
