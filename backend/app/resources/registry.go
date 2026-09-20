package resources

import "goravel/app/resource"

func AdminRegistry() *resource.Registry {
	registry := resource.NewRegistry()
	_ = registry.Register(resource.Manifest{Name: "users", Label: "Users", Route: "/admin/users", Permissions: []string{"admin.users.view"}, Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "email", Label: "Email", Sortable: true}, {Name: "is_active", Label: "Status", Sortable: true}}})
	_ = registry.Register(resource.Manifest{Name: "roles", Label: "Roles", Route: "/admin/roles", Permissions: []string{"admin.roles.manage"}, Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "display_name", Label: "Display name", Sortable: true}}})
	_ = registry.Register(resource.Manifest{Name: "permissions", Label: "Permissions", Route: "/admin/permissions", Permissions: []string{"admin.permissions.manage"}, Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "display_name", Label: "Display name", Sortable: true}}})
	return registry
}
