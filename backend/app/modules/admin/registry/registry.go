package registry

import "goravel/app/core/resource"

func AdminRegistry() *resource.Registry {
	registry := resource.NewRegistry()
	registerGenerated(registry)
	_ = registry.Register(resource.Manifest{
		Name: "users", Label: "Users", Route: "/admin/users", Table: "users", Permissions: []string{"admin.users.view"}, DataScope: resource.DataScopeAll, Navigation: resource.Navigation{Group: "system", Order: 10},
		Fields: []resource.Field{
			{Name: "name", Label: "Name", Type: "text"}, {Name: "email", Label: "Email", Type: "email"},
			{Name: "password", Label: "Password", Type: "password", Visible: false, Readable: false, Writable: true, Sensitive: true, PolicyConfigured: true}, {Name: "locale", Label: "Locale", Type: "text"},
			{Name: "status", Label: "Status", Type: "select", Options: []resource.Option{
				{Value: "active", Label: "Active"}, {Value: "disabled", Label: "Disabled"}, {Value: "locked", Label: "Locked"},
			}},
		},
		Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "email", Label: "Email", Sortable: true}, {Name: "status", Label: "Status", Sortable: true}},
		Actions: append(standardActions("admin.users.manage"), resource.Action{Name: "set-status", Label: "Set status", Kind: "user-status", Permission: "admin.users.manage", Batch: true, Payload: "user-status"}),
		Filters: []resource.Filter{{Name: "status", Label: "Status", Type: "select", Options: []resource.Option{{Value: "active", Label: "Active"}, {Value: "disabled", Label: "Disabled"}, {Value: "locked", Label: "Locked"}}}},
	})
	_ = registry.Register(resource.Manifest{
		Name: "roles", Label: "Roles", Route: "/admin/roles", Table: "roles", Permissions: []string{"admin.roles.manage"}, DataScope: resource.DataScopeAll, Navigation: resource.Navigation{Group: "system", Order: 20},
		Fields:  []resource.Field{{Name: "name", Label: "Name", Type: "text"}, {Name: "display_name", Label: "Display name", Type: "text"}},
		Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "display_name", Label: "Display name", Sortable: true}},
		Actions: standardActions("admin.roles.manage"),
	})
	_ = registry.Register(resource.Manifest{
		Name: "permissions", Label: "Permissions", Route: "/admin/permissions", Table: "permissions", Permissions: []string{"admin.permissions.manage"}, DataScope: resource.DataScopeAll, Navigation: resource.Navigation{Group: "system", Order: 30},
		Fields:  []resource.Field{{Name: "name", Label: "Name", Type: "text"}, {Name: "display_name", Label: "Display name", Type: "text"}},
		Columns: []resource.Column{{Name: "name", Label: "Name", Sortable: true}, {Name: "display_name", Label: "Display name", Sortable: true}},
	})
	return registry
}

func standardActions(permission string) []resource.Action {
	return []resource.Action{
		{Name: "view", Label: "View", Permission: permission},
		{Name: "create", Label: "Create", Permission: permission},
		{Name: "update", Label: "Update", Permission: permission},
		{Name: "delete", Label: "Delete", Permission: permission},
		{Name: "bulk-delete", Label: "Delete selected", Kind: "builtin-delete", Permission: permission, Batch: true, Payload: "delete"},
		{Name: "bulk-update", Label: "Update selected", Kind: "builtin-update", Permission: permission, Batch: true, Payload: "fields"},
	}
}
