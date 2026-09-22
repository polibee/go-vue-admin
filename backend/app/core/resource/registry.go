package resource

import (
	"errors"
	"sort"
)

var (
	ErrInvalidManifest = errors.New("invalid resource manifest")
	ErrDuplicate       = errors.New("resource already registered")
	ErrNotFound        = errors.New("resource not found")
)

type DataScope string

const (
	DataScopeAll DataScope = "all"
	DataScopeOwn DataScope = "own"
)

type Field struct {
	Name             string   `json:"name"`
	Label            string   `json:"label"`
	Type             string   `json:"type"`
	Required         bool     `json:"required,omitempty"`
	Options          []Option `json:"options,omitempty"`
	Visible          bool     `json:"visible"`
	Readable         bool     `json:"readable"`
	Writable         bool     `json:"writable"`
	Sensitive        bool     `json:"sensitive"`
	PolicyConfigured bool     `json:"-"`
}

type Option struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Column struct {
	Name     string `json:"name"`
	Label    string `json:"label"`
	Sortable bool   `json:"sortable"`
}

type Action struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	Kind       string `json:"kind"`
	Permission string `json:"permission"`
}

type Manifest struct {
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	Route       string    `json:"route"`
	Table       string    `json:"table,omitempty"`
	Permissions []string  `json:"permissions"`
	Fields      []Field   `json:"fields"`
	Columns     []Column  `json:"columns"`
	Actions     []Action  `json:"actions,omitempty"`
	DataScope   DataScope `json:"data_scope,omitempty"`
	OwnerField  string    `json:"owner_field,omitempty"`
}

type Registry struct{ manifests map[string]Manifest }

func NewRegistry() *Registry { return &Registry{manifests: make(map[string]Manifest)} }

func (r *Registry) Register(manifest Manifest) error {
	if manifest.Name == "" || manifest.Label == "" || manifest.Route == "" {
		return ErrInvalidManifest
	}
	if manifest.DataScope == "" {
		manifest.DataScope = DataScopeAll
	}
	if manifest.DataScope != DataScopeAll && manifest.DataScope != DataScopeOwn {
		return ErrInvalidManifest
	}
	if manifest.DataScope == DataScopeOwn {
		if manifest.OwnerField == "" {
			return ErrInvalidManifest
		}
		ownerDeclared := false
		for _, field := range manifest.Fields {
			if field.Name == manifest.OwnerField {
				ownerDeclared = true
				break
			}
		}
		if !ownerDeclared {
			return ErrInvalidManifest
		}
	}
	manifest = NormalizeManifestFields(manifest)
	if _, exists := r.manifests[manifest.Name]; exists {
		return ErrDuplicate
	}
	r.manifests[manifest.Name] = manifest
	return nil
}

// NormalizeManifestFields preserves the legacy behavior for fields that do not
// declare a policy while allowing generated or explicit restricted fields to
// retain false values.
func NormalizeManifestFields(manifest Manifest) Manifest {
	for index, field := range manifest.Fields {
		if !field.PolicyConfigured && !field.Visible && !field.Readable && !field.Writable && !field.Sensitive {
			field.Visible = true
			field.Readable = true
			field.Writable = true
		}
		manifest.Fields[index] = field
	}
	return manifest
}

func (r *Registry) Find(name string) (Manifest, error) {
	manifest, ok := r.manifests[name]
	if !ok {
		return Manifest{}, ErrNotFound
	}
	return manifest, nil
}

func (r *Registry) All() []Manifest {
	manifests := make([]Manifest, 0, len(r.manifests))
	for _, manifest := range r.manifests {
		manifests = append(manifests, manifest)
	}
	sort.Slice(manifests, func(i, j int) bool { return manifests[i].Name < manifests[j].Name })
	return manifests
}
