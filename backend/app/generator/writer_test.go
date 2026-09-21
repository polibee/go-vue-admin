package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAllCreatesAllArtifacts(t *testing.T) {
	root := t.TempDir()
	artifacts := []Artifact{
		{Path: "app/generated/resources/posts/manifest.go", Content: []byte("manifest")},
		{Path: "database/migrations/20260921000000_create_posts_table.go", Content: []byte("migration")},
	}

	if err := WriteAll(root, artifacts); err != nil {
		t.Fatalf("write all: %v", err)
	}
	for _, artifact := range artifacts {
		got, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(artifact.Path)))
		if err != nil {
			t.Fatalf("read %s: %v", artifact.Path, err)
		}
		if string(got) != string(artifact.Content) {
			t.Errorf("content for %s = %q", artifact.Path, got)
		}
	}
}

func TestWriteAllDoesNotPartiallyWriteOnConflict(t *testing.T) {
	root := t.TempDir()
	conflict := filepath.Join(root, "app", "generated", "resources", "posts", "model.go")
	if err := os.MkdirAll(filepath.Dir(conflict), 0o755); err != nil {
		t.Fatalf("create conflict directory: %v", err)
	}
	if err := os.WriteFile(conflict, []byte("existing"), 0o644); err != nil {
		t.Fatalf("create conflict: %v", err)
	}

	err := WriteAll(root, []Artifact{
		{Path: "app/generated/resources/posts/manifest.go", Content: []byte("manifest")},
		{Path: "app/generated/resources/posts/model.go", Content: []byte("model")},
	})
	if err == nil {
		t.Fatal("expected conflict error")
	}
	if _, statErr := os.Stat(filepath.Join(root, "app", "generated", "resources", "posts", "manifest.go")); !os.IsNotExist(statErr) {
		t.Fatalf("manifest was partially written, stat error = %v", statErr)
	}
	got, readErr := os.ReadFile(conflict)
	if readErr != nil || string(got) != "existing" {
		t.Fatalf("conflict file changed: %q, %v", got, readErr)
	}
}
