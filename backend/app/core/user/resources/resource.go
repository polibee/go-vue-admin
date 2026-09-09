package resources

const (
	PermissionView   = "users.view"
	PermissionCreate = "users.create"
	PermissionUpdate = "users.update"
	PermissionDelete = "users.delete"
)

type Record struct {
	ID      string   `json:"id"`
	Email   string   `json:"email"`
	Name    string   `json:"name"`
	Active  bool     `json:"active"`
	RoleIDs []string `json:"role_ids"`
}
