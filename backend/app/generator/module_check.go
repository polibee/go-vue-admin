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
	spec, err := Normalize(Input{Name: name, Fields: []string{"name:text"}})
	if err != nil {
		return ModuleCheckReport{}, err
	}
	artifacts, err := RenderResourcePipeline(spec, "00000000000000")
	if err != nil {
		return ModuleCheckReport{}, err
	}
	expected := make([]string, 0, len(artifacts))
	for _, artifact := range artifacts {
		path := filepath.ToSlash(artifact.Path)
		if path == "database/migrations/00000000000000_create_"+name+"_table.go" {
			continue
		}
		expected = append(expected, path)
	}
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
