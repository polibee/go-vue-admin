package generator

import (
	"os"
	"path/filepath"
	"strings"
)

func GenerateResource(root string, input Input, migrationTimestamp string) ([]Artifact, error) {
	spec, err := Normalize(input)
	if err != nil {
		return nil, err
	}
	artifacts, err := RenderResourcePipeline(spec, migrationTimestamp)
	if err != nil {
		return nil, err
	}
	runtimeArtifacts, err := RenderRuntimeRegistration(root, spec)
	if err != nil {
		return nil, err
	}
	artifacts = append(artifacts, runtimeArtifacts...)
	if err := writeProjectArtifacts(root, artifacts); err != nil {
		return nil, err
	}
	return artifacts, nil
}

func writeProjectArtifacts(root string, artifacts []Artifact) error {
	frontendRoot := resolveFrontendRoot(root)
	if frontendRoot == root {
		return WriteAll(root, artifacts)
	}

	backendArtifacts := make([]Artifact, 0, len(artifacts))
	frontendArtifacts := make([]Artifact, 0)
	for _, artifact := range artifacts {
		path := filepath.ToSlash(artifact.Path)
		if strings.HasPrefix(path, "admin/") {
			frontendArtifacts = append(frontendArtifacts, Artifact{
				Path:           strings.TrimPrefix(path, "admin/"),
				Content:        artifact.Content,
				AllowOverwrite: artifact.AllowOverwrite,
			})
			continue
		}
		backendArtifacts = append(backendArtifacts, artifact)
	}
	if err := WriteAll(root, backendArtifacts); err != nil {
		return err
	}
	return WriteAll(frontendRoot, frontendArtifacts)
}

func resolveFrontendRoot(backendRoot string) string {
	cleanRoot := filepath.Clean(backendRoot)
	if filepath.Base(cleanRoot) != "backend" {
		return backendRoot
	}
	frontendRoot := filepath.Join(filepath.Dir(cleanRoot), "admin")
	if info, err := os.Stat(frontendRoot); err == nil && info.IsDir() {
		return frontendRoot
	}
	return backendRoot
}
