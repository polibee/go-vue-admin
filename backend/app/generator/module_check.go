package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

type ModuleCheckReport struct {
	Name     string
	Expected []string
	Missing  []string
	Complete bool
}

func CheckModule(root, name string) (ModuleCheckReport, error) {
	spec, err := NormalizeModule(name)
	if err != nil {
		return ModuleCheckReport{}, err
	}
	expected := moduleArtifactPaths(spec)
	report := ModuleCheckReport{Name: spec.Name, Expected: expected}
	for _, relative := range expected {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(relative))); err != nil {
			if os.IsNotExist(err) {
				report.Missing = append(report.Missing, relative)
				continue
			}
			return ModuleCheckReport{}, fmt.Errorf("check %s: %w", relative, err)
		}
	}
	report.Complete = len(report.Missing) == 0
	return report, nil
}

func moduleArtifactPaths(spec ModuleSpec) []string {
	base := filepath.ToSlash(filepath.Join("app", "modules", spec.Name))
	paths := []string{"module.go", "model.go", "request.go", "repository.go", "service.go", "controller.go", "routes.go", "resource.go", "permissions.go", "events.go", "tests/module_test.go", "README.md"}
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		result = append(result, filepath.ToSlash(filepath.Join(base, path)))
	}
	return result
}
