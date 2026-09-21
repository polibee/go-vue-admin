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
	if len(command.Extend().Flags) != 4 {
		t.Fatalf("flag count = %d, want 4", len(command.Extend().Flags))
	}
}
