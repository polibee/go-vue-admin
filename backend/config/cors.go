package config

import (
	"goravel/app/facades"
)

func init() {
	config := facades.Config()
	config.Add("cors", map[string]any{
		// Cross-Origin Resource Sharing (CORS) Configuration
		//
		// Here you may configure your settings for cross-origin resource sharing
		// or "CORS". This determines what cross-origin operations may execute
		// in web browsers. You are free to adjust these settings as needed.
		//
		// To learn more: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
		"paths":                []string{"api/*"},
		"allowed_methods":      []string{"*"},
		"allowed_origins":      []string{"http://127.0.0.1:5187", "http://localhost:5187", "http://127.0.0.1:5173", "http://localhost:5173"},
		"allowed_headers":      []string{"*"},
		"exposed_headers":      []string{},
		"max_age":              0,
		"supports_credentials": true,
	})
}
