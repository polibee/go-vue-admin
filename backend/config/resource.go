package config

import "goravel/app/facades"

func init() {
	config := facades.Config()
	config.Add("resource", map[string]any{
		"provider": config.Env("RESOURCE_PROVIDER", "memory"),
	})
}
