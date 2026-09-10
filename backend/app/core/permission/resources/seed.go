package resources

import "goravel/app/core/resource"

func Seed() []resource.CoreResourceRecord {
	return []resource.CoreResourceRecord{
		{"id": "dashboard.view", "name": "dashboard.view", "description": "查看仪表盘"},
		{"id": "settings.view", "name": "settings.view", "description": "查看设置"},
		{"id": "settings.update", "name": "settings.update", "description": "修改设置"},
		{"id": "users.view", "name": "users.view", "description": "查看用户"},
		{"id": "users.create", "name": "users.create", "description": "创建用户"},
		{"id": "users.update", "name": "users.update", "description": "编辑用户"},
		{"id": "users.delete", "name": "users.delete", "description": "删除用户"},
		{"id": "roles.view", "name": "roles.view", "description": "查看角色"},
		{"id": "roles.create", "name": "roles.create", "description": "创建角色"},
		{"id": "roles.update", "name": "roles.update", "description": "编辑角色"},
		{"id": "roles.delete", "name": "roles.delete", "description": "删除角色"},
		{"id": "permissions.view", "name": "permissions.view", "description": "查看权限"},
		{"id": "permissions.create", "name": "permissions.create", "description": "创建权限"},
		{"id": "permissions.update", "name": "permissions.update", "description": "编辑权限"},
		{"id": "permissions.delete", "name": "permissions.delete", "description": "删除权限"},
	}
}
