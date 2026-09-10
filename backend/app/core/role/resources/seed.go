package resources

import "goravel/app/core/resource"

func Seed() []resource.CoreResourceRecord {
	return []resource.CoreResourceRecord{
		{
			"id": "platform-admin", "name": "平台管理员", "description": "访问平台资源的开发管理员",
			"permissions": []string{"dashboard.view", "settings.view", "settings.update", "media.*", "audit.view", "users.*", "roles.*", "permissions.*"},
		},
	}
}
