package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"goravel/app/generator"
)

type MenuGeneratorCommand struct{}

func (MenuGeneratorCommand) Signature() string { return "admin:make-menu" }

func (MenuGeneratorCommand) Description() string {
	return "Generate a reviewable menu entry"
}

func (MenuGeneratorCommand) Extend() command.Extend {
	return command.Extend{
		Category:  "admin",
		ArgsUsage: "<name>",
		Arguments: []command.Argument{&command.ArgumentString{Name: "name", Usage: "menu name", Required: true}},
		Flags: []command.Flag{
			&command.StringFlag{Name: "label", Usage: "menu label", Required: true},
			&command.StringFlag{Name: "route", Usage: "menu route", Required: true},
			&command.StringFlag{Name: "permission", Usage: "required permission", Required: true},
			&command.StringFlag{Name: "icon", Usage: "menu icon", Required: true},
		},
	}
}

func (MenuGeneratorCommand) Handle(ctx console.Context) error {
	spec, err := generator.NormalizeMenu(generator.MenuInput{
		Name:       ctx.ArgumentString("name"),
		Label:      ctx.Option("label"),
		Route:      ctx.Option("route"),
		Permission: ctx.Option("permission"),
		Icon:       ctx.Option("icon"),
	})
	if err != nil {
		return err
	}
	artifacts, err := generator.RenderMenu(spec)
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
	ctx.Success("Generated menu skeleton:")
	for _, artifact := range artifacts {
		ctx.Line("- " + strings.ReplaceAll(artifact.Path, "\\", "/"))
	}
	return nil
}
