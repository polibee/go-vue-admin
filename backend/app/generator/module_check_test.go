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
	spec, err := NormalizeModule("billing")
	if err != nil {
		t.Fatal(err)
	}
	artifacts, err := RenderModule(spec)
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
	if !report.Complete || len(report.Missing) != 0 {
		t.Fatalf("report = %+v", report)
	}
}
