package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"goravel/app/generator"
)

type ModuleGeneratorCommand struct{}

func (ModuleGeneratorCommand) Signature() string { return "admin:make-module" }

func (ModuleGeneratorCommand) Description() string {
	return "Generate a reviewable backend module skeleton"
}

func (ModuleGeneratorCommand) Extend() command.Extend {
	return command.Extend{
		Category:  "admin",
		ArgsUsage: "<name>",
		Arguments: []command.Argument{&command.ArgumentString{Name: "name", Usage: "module name", Required: true}},
	}
}

func (ModuleGeneratorCommand) Handle(ctx console.Context) error {
	spec, err := generator.NormalizeModule(ctx.ArgumentString("name"))
	if err != nil {
		return err
	}
	artifacts, err := generator.RenderModule(spec)
	if err != nil {
		return err
	}
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	if err := generator.WriteAll(root, artifacts); err != nil {
		return err
	}
	ctx.Success("Generated module skeleton:")
	for _, artifact := range artifacts {
		ctx.Line("- " + strings.ReplaceAll(artifact.Path, "\\", "/"))
	}
	return nil
}
