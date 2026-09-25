package generator

func RenderResourcePipeline(spec Spec, migrationTimestamp string) ([]Artifact, error) {
	resourceArtifacts, err := Render(spec, migrationTimestamp)
	if err != nil {
		return nil, err
	}
	permissionSpec, err := NormalizePermission(PermissionInput{Name: spec.Name, Namespace: spec.Namespace, Actions: spec.Actions})
	if err != nil {
		return nil, err
	}
	permissionArtifacts, err := RenderPermission(permissionSpec)
	if err != nil {
		return nil, err
	}
	menuSpec, err := NormalizeMenu(MenuInput{
		Name:       spec.Name,
		Namespace:  spec.Namespace,
		Label:      spec.Label,
		Route:      spec.FrontendRoute,
		Permission: spec.Permission,
		Icon:       spec.Icon,
	})
	if err != nil {
		return nil, err
	}
	menuArtifacts, err := RenderMenu(menuSpec)
	if err != nil {
		return nil, err
	}
	frontendArtifacts, err := RenderFrontend(spec)
	if err != nil {
		return nil, err
	}
	artifacts := append(resourceArtifacts, permissionArtifacts...)
	artifacts = append(artifacts, menuArtifacts...)
	return append(artifacts, frontendArtifacts...), nil
}
