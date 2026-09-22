package console

import "testing"

func TestResourceGeneratorCommandMetadata(t *testing.T) {
	command := ResourceGeneratorCommand{}
	if command.Signature() != "admin:make-resource" {
		t.Fatalf("signature = %q", command.Signature())
	}
	if command.Description() == "" {
		t.Fatal("description must not be empty")
	}
	if len(command.Extend().Flags) != 10 {
		t.Fatalf("flag count = %d, want 10", len(command.Extend().Flags))
	}
}

func TestModuleCheckCommandMetadata(t *testing.T) {
	command := ModuleCheckCommand{}
	if command.Signature() != "admin:check-module" {
		t.Fatalf("signature = %q", command.Signature())
	}
	if len(command.Extend().Arguments) != 1 {
		t.Fatalf("argument count = %d, want 1", len(command.Extend().Arguments))
	}
}
