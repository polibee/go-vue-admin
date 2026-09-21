package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRenderPostsGolden(t *testing.T) {
	spec, err := Normalize(Input{
		Name:   "posts",
		Label:  "Posts",
		Route:  "/admin/posts",
		Fields: []string{"title:text:required", "published:boolean"},
	})
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}

	artifacts, err := Render(spec, "20260921000000")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	wantPaths := []string{
		"app/generated/resources/posts/manifest.go",
		"app/generated/resources/posts/model.go",
		"app/generated/resources/posts/request.go",
		"app/generated/resources/posts/repository.go",
		"app/generated/resources/posts/service.go",
		"app/generated/resources/posts/controller.go",
		"app/generated/resources/posts/routes.go",
		"app/generated/resources/posts/permissions.go",
		"app/generated/resources/posts/manifest_test.go",
		"database/migrations/20260921000000_create_posts_table.go",
	}
	if got := artifactPaths(artifacts); !reflect.DeepEqual(got, wantPaths) {
		t.Fatalf("artifact paths = %v, want %v", got, wantPaths)
	}

	for _, artifact := range artifacts {
		goldenPath := filepath.Join("testdata", "posts", filepath.Base(artifact.Path)+".golden")
		if os.Getenv("UPDATE_GOLDEN") == "1" {
			if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
				t.Fatalf("create golden directory: %v", err)
			}
			if err := os.WriteFile(goldenPath, artifact.Content, 0o644); err != nil {
				t.Fatalf("write golden %s: %v", goldenPath, err)
			}
			continue
		}
		want, readErr := os.ReadFile(goldenPath)
		if readErr != nil {
			t.Fatalf("read golden %s: %v", goldenPath, readErr)
		}
		if string(artifact.Content) != string(want) {
			t.Errorf("content mismatch for %s", artifact.Path)
		}
	}

	repeated, err := Render(spec, "20260921000000")
	if err != nil {
		t.Fatalf("render repeated: %v", err)
	}
	if !reflect.DeepEqual(artifacts, repeated) {
		t.Fatal("render output is not deterministic")
	}
}

func artifactPaths(artifacts []Artifact) []string {
	paths := make([]string, 0, len(artifacts))
	for _, artifact := range artifacts {
		paths = append(paths, artifact.Path)
	}
	return paths
}
