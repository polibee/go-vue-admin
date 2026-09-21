package console

import "testing"

func TestPermissionGeneratorCommandMetadata(t *testing.T) {
	command := PermissionGeneratorCommand{}
	if command.Signature() != "admin:make-permission" {
		t.Fatalf("signature = %q", command.Signature())
	}
	if len(command.Extend().Flags) != 1 || len(command.Extend().Arguments) != 1 {
		t.Fatalf("unexpected command shape: %+v", command.Extend())
	}
}
