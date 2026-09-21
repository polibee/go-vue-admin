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

type Field struct {
	Name    string   `json:"name"`
	Label   string   `json:"label"`
	Type    string   `json:"type"`
	Options []Option `json:"options,omitempty"`
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

type Manifest struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Route       string   `json:"route"`
	Permissions []string `json:"permissions"`
	Fields      []Field  `json:"fields"`
	Columns     []Column `json:"columns"`
}

type Registry struct{ manifests map[string]Manifest }

func NewRegistry() *Registry { return &Registry{manifests: make(map[string]Manifest)} }

func (r *Registry) Register(manifest Manifest) error {
	if manifest.Name == "" || manifest.Label == "" || manifest.Route == "" {
		return ErrInvalidManifest
	}
	if _, exists := r.manifests[manifest.Name]; exists {
		return ErrDuplicate
	}
	r.manifests[manifest.Name] = manifest
	return nil
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
