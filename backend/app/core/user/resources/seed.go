package resources

import "goravel/app/core/resource"

func Seed() []resource.CoreResourceRecord {
	return []resource.CoreResourceRecord{
		{
			"id": "bootstrap-admin", "email": "admin@example.com", "name": "Platform Admin",
			"active": true, "role_ids": []string{"platform-admin"},
		},
	}
}
