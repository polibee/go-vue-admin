package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
		"app/modules/posts/resource/manifest.go",
		"app/modules/posts/resource/model.go",
		"app/modules/posts/resource/request.go",
		"app/modules/posts/resource/repository.go",
		"app/modules/posts/resource/service.go",
		"app/modules/posts/resource/controller.go",
		"app/modules/posts/resource/routes.go",
		"app/modules/posts/resource/permissions.go",
		"app/modules/posts/resource/manifest_test.go",
		"database/migrations/20260921000000_create_posts_table.go",
		"app/modules/posts/resource/README.md",
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

func TestRenderActionPayloadFields(t *testing.T) {
	spec, err := Normalize(Input{
		Name:             "orders",
		Fields:           []string{"number:text"},
		ActionSpecs:      []string{"archive:Archive:archive:admin.orders.archive:true:archive"},
		ActionFieldSpecs: []string{"archive:reason:Reason:text:true", "archive:mode:Mode:select:false:fast=Fast|safe=Safe"},
	})
	if err != nil {
		t.Fatalf("normalize: %v", err)
	}
	artifacts, err := Render(spec, "20260922000000")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	for _, artifact := range artifacts {
		if artifact.Path != "app/modules/orders/resource/manifest.go" {
			continue
		}
		content := string(artifact.Content)
		for _, fragment := range []string{"PayloadFields", `Name: "reason"`, `Name: "mode"`, `Value: "fast"`} {
			if !strings.Contains(content, fragment) {
				t.Fatalf("manifest payload field missing %q: %s", fragment, content)
			}
		}
		return
	}
	t.Fatal("manifest artifact not found")
}

func artifactPaths(artifacts []Artifact) []string {
	paths := make([]string, 0, len(artifacts))
	for _, artifact := range artifacts {
		paths = append(paths, artifact.Path)
	}
	return paths
}
