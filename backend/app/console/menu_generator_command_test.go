package console

import "testing"

func TestMenuGeneratorCommandMetadata(t *testing.T) {
	command := MenuGeneratorCommand{}
	if command.Signature() != "admin:make-menu" {
		t.Fatalf("signature = %q", command.Signature())
	}
	if len(command.Extend().Flags) != 4 || len(command.Extend().Arguments) != 1 {
		t.Fatalf("unexpected command shape: %+v", command.Extend())
	}
}
