package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteAll(root string, artifacts []Artifact) error {
	if len(artifacts) == 0 {
		return fmt.Errorf("no artifacts to write")
	}
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		if err == nil {
			err = fmt.Errorf("not a directory")
		}
		return fmt.Errorf("invalid output root: %w", err)
	}

	seen := make(map[string]struct{}, len(artifacts))
	for _, artifact := range artifacts {
		relative, err := safeRelativePath(artifact.Path)
		if err != nil {
			return err
		}
		if _, exists := seen[relative]; exists {
			return fmt.Errorf("duplicate artifact path %q", artifact.Path)
		}
		seen[relative] = struct{}{}
		target := filepath.Join(root, filepath.FromSlash(relative))
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("refusing to overwrite existing file %q", artifact.Path)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check artifact %q: %w", artifact.Path, err)
		}
	}

	temporary, err := os.MkdirTemp(root, ".admin-generator-")
	if err != nil {
		return fmt.Errorf("create temporary output: %w", err)
	}
	defer os.RemoveAll(temporary)

	for _, artifact := range artifacts {
		relative, _ := safeRelativePath(artifact.Path)
		target := filepath.Join(temporary, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create directory for %q: %w", artifact.Path, err)
		}
		if err := os.WriteFile(target, artifact.Content, 0o644); err != nil {
			return fmt.Errorf("write artifact %q: %w", artifact.Path, err)
		}
	}

	for _, artifact := range artifacts {
		relative, _ := safeRelativePath(artifact.Path)
		source := filepath.Join(temporary, filepath.FromSlash(relative))
		target := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create target directory for %q: %w", artifact.Path, err)
		}
		if err := os.Rename(source, target); err != nil {
			return fmt.Errorf("install artifact %q: %w", artifact.Path, err)
		}
	}
	return nil
}

func safeRelativePath(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("artifact path cannot be empty")
	}
	normalized := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if normalized == "." || filepath.IsAbs(filepath.FromSlash(normalized)) || normalized == ".." || strings.HasPrefix(normalized, "../") {
		return "", fmt.Errorf("unsafe artifact path %q", value)
	}
	return normalized, nil
}
