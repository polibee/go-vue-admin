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

type Filter struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Options  []Option `json:"options,omitempty"`
	Relation string   `json:"relation,omitempty"`
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

type Relation struct {
	Name         string   `json:"name"`
	Kind         string   `json:"kind"`
	Resource     string   `json:"resource"`
	Field        string   `json:"field"`
	ForeignField string   `json:"foreign_field"`
	LabelField   string   `json:"label_field"`
	Selectable   bool     `json:"selectable"`
	Multiple     bool     `json:"multiple"`
	Permission   string   `json:"permission,omitempty"`
	FilterFields []string `json:"filter_fields,omitempty"`
}

type FormGroup struct {
	Name    string   `json:"name"`
	Label   string   `json:"label"`
	Columns int      `json:"columns,omitempty"`
	Fields  []string `json:"fields"`
}

type DetailSection struct {
	Name   string   `json:"name"`
	Label  string   `json:"label"`
	Fields []string `json:"fields"`
}

type FieldDependency struct {
	Field string `json:"field"`
	On    string `json:"on"`
	Value string `json:"value"`
}

type Action struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	Kind       string `json:"kind"`
	Permission string `json:"permission"`
	Batch      bool   `json:"batch"`
	Payload    string `json:"payload,omitempty"`
}

type Navigation struct {
	Group  string `json:"group"`
	Order  int    `json:"order"`
	Hidden bool   `json:"hidden,omitempty"`
}

type Manifest struct {
	Name         string            `json:"name"`
	Label        string            `json:"label"`
	Route        string            `json:"route"`
	Table        string            `json:"table,omitempty"`
	Permissions  []string          `json:"permissions"`
	Fields       []Field           `json:"fields"`
	Columns      []Column          `json:"columns"`
	Actions      []Action          `json:"actions,omitempty"`
	Filters      []Filter          `json:"filters,omitempty"`
	Relations    []Relation        `json:"relations,omitempty"`
	FormGroups   []FormGroup       `json:"form_groups,omitempty"`
	Details      []DetailSection   `json:"details,omitempty"`
	Dependencies []FieldDependency `json:"dependencies,omitempty"`
	DataScope    DataScope         `json:"data_scope,omitempty"`
	OwnerField   string            `json:"owner_field,omitempty"`
	SoftDelete   bool              `json:"soft_delete,omitempty"`
	Navigation   Navigation        `json:"navigation"`
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
	if manifest.Navigation.Group == "" {
		manifest.Navigation.Group = "business"
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
	if err := ValidateManifestExtensions(manifest); err != nil {
		return err
	}
	if _, exists := r.manifests[manifest.Name]; exists {
		return ErrDuplicate
	}
	r.manifests[manifest.Name] = manifest
	return nil
}

func ValidateManifestExtensions(manifest Manifest) error {
	fields := make(map[string]struct{}, len(manifest.Fields))
	for _, field := range manifest.Fields {
		fields[field.Name] = struct{}{}
	}
	filterNames := make(map[string]struct{}, len(manifest.Filters))
	for _, filter := range manifest.Filters {
		if filter.Name == "" || filter.Label == "" {
			return ErrInvalidManifest
		}
		if filter.Type != "select" && filter.Type != "multi-select" && filter.Type != "boolean" && filter.Type != "text" && filter.Type != "date-range" && filter.Type != "relation" {
			return ErrInvalidManifest
		}
		if _, exists := fields[filter.Name]; !exists && filter.Type != "relation" {
			return ErrInvalidManifest
		}
		if _, exists := filterNames[filter.Name]; exists {
			return ErrInvalidManifest
		}
		filterNames[filter.Name] = struct{}{}
	}
	relations := make(map[string]struct{}, len(manifest.Relations))
	for _, relation := range manifest.Relations {
		if relation.Name == "" || relation.Resource == "" || relation.Field == "" || relation.ForeignField == "" || relation.LabelField == "" {
			return ErrInvalidManifest
		}
		if relation.Kind != "belongsTo" && relation.Kind != "hasMany" {
			return ErrInvalidManifest
		}
		if _, exists := relations[relation.Name]; exists {
			return ErrInvalidManifest
		}
		relations[relation.Name] = struct{}{}
		if relation.Kind == "belongsTo" {
			if _, exists := fields[relation.Field]; !exists || relation.Multiple {
				return ErrInvalidManifest
			}
		} else if relation.Selectable || relation.Multiple {
			return ErrInvalidManifest
		}
	}
	for _, group := range manifest.FormGroups {
		if group.Name == "" || group.Label == "" || group.Columns < 0 || group.Columns > 4 || len(group.Fields) == 0 {
			return ErrInvalidManifest
		}
		if err := validateFieldReferences(group.Fields, fields); err != nil {
			return err
		}
	}
	for _, section := range manifest.Details {
		if section.Name == "" || section.Label == "" || len(section.Fields) == 0 {
			return ErrInvalidManifest
		}
		if err := validateFieldReferences(section.Fields, fields); err != nil {
			return err
		}
	}
	for _, dependency := range manifest.Dependencies {
		if dependency.Field == "" || dependency.On == "" || dependency.Value == "" {
			return ErrInvalidManifest
		}
		if _, exists := fields[dependency.Field]; !exists {
			return ErrInvalidManifest
		}
		if _, exists := fields[dependency.On]; !exists {
			return ErrInvalidManifest
		}
	}
	return nil
}

func validateFieldReferences(names []string, fields map[string]struct{}) error {
	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		if _, exists := fields[name]; !exists {
			return ErrInvalidManifest
		}
		if _, exists := seen[name]; exists {
			return ErrInvalidManifest
		}
		seen[name] = struct{}{}
	}
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
