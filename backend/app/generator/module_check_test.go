package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckModuleReportsMissingFilesWithoutWriting(t *testing.T) {
	root := t.TempDir()
	report, err := CheckModule(root, "billing")
	if err != nil {
		t.Fatalf("check module: %v", err)
	}
	if report.Complete {
		t.Fatal("expected incomplete module report")
	}
	if len(report.Missing) == 0 {
		t.Fatal("expected missing files")
	}
	entries, err := os.ReadDir(filepath.Join(root, "app", "modules"))
	if err == nil || entries != nil {
		t.Fatalf("check command wrote files: entries=%v err=%v", entries, err)
	}
}

func TestCheckModuleReportsCompleteGeneratedModule(t *testing.T) {
	root := t.TempDir()
	spec, err := Normalize(Input{Name: "billing", Fields: []string{"name:text"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderResourcePipeline(spec, "00000000000000")
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteAll(root, artifacts); err != nil {
		t.Fatal(err)
	}
	runtimeArtifacts, err := RenderRuntimeRegistration(root, spec)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteAll(root, runtimeArtifacts); err != nil {
		t.Fatal(err)
	}
	report, err := CheckModule(root, "billing")
	if err != nil {
		t.Fatalf("check module: %v", err)
	}
	if !report.Complete || len(report.Missing) != 0 {
		t.Fatalf("report = %+v", report)
	}
}

func TestCheckModuleRequiresRuntimeDiscoveryFiles(t *testing.T) {
	root := t.TempDir()
	spec, err := Normalize(Input{Name: "billing", Fields: []string{"name:text"}})
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderResourcePipeline(spec, "00000000000000")
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteAll(root, artifacts); err != nil {
		t.Fatal(err)
	}
	report, err := CheckModule(root, "billing")
	if err != nil {
		t.Fatalf("check module: %v", err)
	}
	if report.Complete || !containsPath(report.Missing, "app/modules/admin/registry/generated_resources.go") {
		t.Fatalf("report = %+v, want missing backend discovery", report)
	}
}

func containsPath(paths []string, want string) bool {
	for _, path := range paths {
		if path == want {
			return true
		}
	}
	return false
}
