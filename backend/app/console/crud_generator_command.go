package console

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"goravel/app/generator"
)

type CRUDGeneratorCommand struct{}

func (CRUDGeneratorCommand) Signature() string { return "admin:make-crud" }

func (CRUDGeneratorCommand) Description() string {
	return "Generate a reviewable backend module and resource CRUD skeleton"
}

func (CRUDGeneratorCommand) Extend() command.Extend {
	return command.Extend{
		Category:  "admin",
		ArgsUsage: "<name>",
		Arguments: []command.Argument{&command.ArgumentString{Name: "name", Usage: "CRUD name", Required: true}},
		Flags: []command.Flag{
			&command.StringFlag{Name: "label", Usage: "resource label"},
			&command.StringFlag{Name: "route", Usage: "admin route"},
			&command.StringFlag{Name: "permission", Usage: "view permission"},
			&command.StringSliceFlag{Name: "field", Usage: "field definition name:type[:required]", Required: true},
		},
	}
}

func (CRUDGeneratorCommand) Handle(ctx console.Context) error {
	spec, err := generator.NormalizeCRUD(generator.Input{
		Name:       ctx.ArgumentString("name"),
		Label:      ctx.Option("label"),
		Route:      ctx.Option("route"),
		Permission: ctx.Option("permission"),
		Fields:     ctx.OptionSlice("field"),
	})
	if err != nil {
		return err
	}
	artifacts, err := generator.RenderCRUD(spec, time.Now().UTC().Format("20060102150405"))
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
	ctx.Success("Generated CRUD skeleton:")
	for _, artifact := range artifacts {
		ctx.Line("- " + strings.ReplaceAll(artifact.Path, "\\", "/"))
	}
	return nil
}
