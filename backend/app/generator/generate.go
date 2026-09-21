package generator

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
	if err := WriteAll(root, artifacts); err != nil {
		return nil, err
	}
	return artifacts, nil
}
