package generator

type CRUDSpec struct {
	Resource Spec
	Module   ModuleSpec
}

func NormalizeCRUD(input Input) (CRUDSpec, error) {
	resource, err := Normalize(input)
	if err != nil {
		return CRUDSpec{}, err
	}
	module, err := NormalizeModule(input.Name)
	if err != nil {
		return CRUDSpec{}, err
	}
	return CRUDSpec{Resource: resource, Module: module}, nil
}

func RenderCRUD(spec CRUDSpec, migrationTimestamp string) ([]Artifact, error) {
	moduleArtifacts, err := RenderModule(spec.Module)
	if err != nil {
		return nil, err
	}
	resourceArtifacts, err := Render(spec.Resource, migrationTimestamp)
	if err != nil {
		return nil, err
	}
	return append(moduleArtifacts, resourceArtifacts...), nil
}
