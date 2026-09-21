package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"goravel/app/generator"
)

type PermissionGeneratorCommand struct{}

func (PermissionGeneratorCommand) Signature() string { return "admin:make-permission" }

func (PermissionGeneratorCommand) Description() string {
	return "Generate reviewable permission constants"
}

func (PermissionGeneratorCommand) Extend() command.Extend {
	return command.Extend{
		Category:  "admin",
		ArgsUsage: "<name>",
		Arguments: []command.Argument{&command.ArgumentString{Name: "name", Usage: "permission resource name", Required: true}},
		Flags: []command.Flag{
			&command.StringSliceFlag{Name: "action", Usage: "permission action", Required: true},
		},
	}
}

func (PermissionGeneratorCommand) Handle(ctx console.Context) error {
	spec, err := generator.NormalizePermission(generator.PermissionInput{Name: ctx.ArgumentString("name"), Actions: ctx.OptionSlice("action")})
	if err != nil {
		return err
	}
	artifacts, err := generator.RenderPermission(spec)
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
	ctx.Success("Generated permission skeleton:")
	for _, artifact := range artifacts {
		ctx.Line("- " + strings.ReplaceAll(artifact.Path, "\\", "/"))
	}
	return nil
}
