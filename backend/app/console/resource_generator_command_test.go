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
	if len(command.Extend().Flags) != 5 {
		t.Fatalf("flag count = %d, want 5", len(command.Extend().Flags))
	}
}

func TestModuleGeneratorCommandMetadata(t *testing.T) {
	command := ModuleGeneratorCommand{}
	if command.Signature() != "admin:make-module" {
		t.Fatalf("signature = %q", command.Signature())
	}
	if len(command.Extend().Arguments) != 1 {
		t.Fatalf("argument count = %d, want 1", len(command.Extend().Arguments))
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
