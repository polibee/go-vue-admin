package resources

const (
	PermissionView   = "roles.view"
	PermissionCreate = "roles.create"
	PermissionUpdate = "roles.update"
	PermissionDelete = "roles.delete"
)

type Record struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}
