package generator

import "fmt"

type ModuleSpec struct {
	Name   string
	GoName string
}

func NormalizeModule(name string) (ModuleSpec, error) {
	if !resourceNamePattern.MatchString(name) {
		return ModuleSpec{}, fmt.Errorf("invalid module name %q", name)
	}
	return ModuleSpec{Name: name, GoName: pascal(name)}, nil
}
