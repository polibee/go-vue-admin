package console

import (
	"testing"

	"github.com/goravel/framework/contracts/console/command"
)

func TestResourceGeneratorCommandMetadata(t *testing.T) {
	generatorCommand := ResourceGeneratorCommand{}
	if generatorCommand.Signature() != "admin:make-resource" {
		t.Fatalf("signature = %q", generatorCommand.Signature())
	}
	if generatorCommand.Description() == "" {
		t.Fatal("description must not be empty")
	}
	if len(generatorCommand.Extend().Flags) != 14 {
		t.Fatalf("flag count = %d, want 14", len(generatorCommand.Extend().Flags))
	}
	flag, ok := generatorCommand.Extend().Flags[0].(*command.StringFlag)
	if !ok || flag.Name != "namespace" {
		t.Fatalf("first flag is not namespace")
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
