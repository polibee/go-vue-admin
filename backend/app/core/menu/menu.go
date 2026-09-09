package menu

import "goravel/app/core/permission"

type Item struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Route      string `json:"route"`
	Permission string `json:"permission,omitempty"`
}

type Registry struct {
	items []Item
}

func NewRegistry(items ...Item) *Registry {
	registry := &Registry{items: make([]Item, len(items))}
	copy(registry.items, items)
	return registry
}

func (r *Registry) All() []Item {
	items := make([]Item, len(r.items))
	copy(items, r.items)
	return items
}

func (r *Registry) ForPermissions(granted []string) []Item {
	authorizer := permission.NewAuthorizer()
	items := make([]Item, 0, len(r.items))
	for _, item := range r.items {
		if item.Permission == "" || authorizer.Allows(granted, item.Permission) {
			items = append(items, item)
		}
	}
	return items
}
