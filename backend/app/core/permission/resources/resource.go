package resources

const (
	PermissionView   = "permissions.view"
	PermissionCreate = "permissions.create"
	PermissionUpdate = "permissions.update"
	PermissionDelete = "permissions.delete"
)

type Record struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
