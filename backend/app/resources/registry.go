package resources

import "goravel/app/resource"

func AdminRegistry() *resource.Registry {
	registry := resource.NewRegistry()
	_ = registry.Register(resource.Manifest{
		Name: "users", Label: "Users", Route: "/admin/users", Permissions: []string{"admin.users.view"},
		Fields: []resource.Field{
			{Name: "name", Label: "Name", Type: "text"}, {Name: "email", Label: "Email", Type: "email"},
			{Name: "password", Label: "Password", Type: "password"}, {Name: "locale", Label: "Locale", Type: "text"},
			{Name: "status", Label: "Status", Type: "select", Options: []resource.Option{
				{Value: "active", Label: "Active"}, {Value: "disabled", Label: "Disabled"}, {Value: "locked", Label: "Locked"},
			}},
		},
		Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "email", Label: "Email", Sortable: true}, {Name: "status", Label: "Status", Sortable: true}},
	})
	_ = registry.Register(resource.Manifest{
		Name: "roles", Label: "Roles", Route: "/admin/roles", Permissions: []string{"admin.roles.manage"},
		Fields:  []resource.Field{{Name: "name", Label: "Name", Type: "text"}, {Name: "display_name", Label: "Display name", Type: "text"}},
		Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "display_name", Label: "Display name", Sortable: true}},
	})
	_ = registry.Register(resource.Manifest{
		Name: "permissions", Label: "Permissions", Route: "/admin/permissions", Permissions: []string{"admin.permissions.manage"},
		Fields:  []resource.Field{{Name: "name", Label: "Name", Type: "text"}, {Name: "display_name", Label: "Display name", Type: "text"}},
		Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "display_name", Label: "Display name", Sortable: true}},
	})
	return registry
}
