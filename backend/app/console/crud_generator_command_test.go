package console

import "testing"

func TestCRUDGeneratorCommandMetadata(t *testing.T) {
	command := CRUDGeneratorCommand{}
	if command.Signature() != "admin:make-crud" {
		t.Fatalf("signature = %q", command.Signature())
	}
	if len(command.Extend().Flags) != 5 || len(command.Extend().Arguments) != 1 {
		t.Fatalf("unexpected command shape: %+v", command.Extend())
	}
}
