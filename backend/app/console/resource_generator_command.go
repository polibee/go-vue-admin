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

type ResourceGeneratorCommand struct{}

func (ResourceGeneratorCommand) Signature() string {
	return "admin:make-resource"
}

func (ResourceGeneratorCommand) Description() string {
	return "Generate the complete resource-driven backend and frontend skeleton"
}

func (ResourceGeneratorCommand) Extend() command.Extend {
	return command.Extend{
		Category:  "admin",
		ArgsUsage: "<name>",
		Arguments: []command.Argument{&command.ArgumentString{Name: "name", Usage: "resource name", Required: true}},
		Flags: []command.Flag{
			&command.StringFlag{Name: "label", Usage: "resource label"},
			&command.StringFlag{Name: "route", Usage: "admin route"},
			&command.StringFlag{Name: "permission", Usage: "view permission"},
			&command.StringFlag{Name: "icon", Usage: "menu icon"},
			&command.StringSliceFlag{Name: "field", Usage: "field definition name:type[:required[:value=Label|value=Label]]", Required: true},
		},
	}
}

func (ResourceGeneratorCommand) Handle(ctx console.Context) error {
	input := generator.Input{
		Name:       ctx.ArgumentString("name"),
		Label:      ctx.Option("label"),
		Route:      ctx.Option("route"),
		Permission: ctx.Option("permission"),
		Icon:       ctx.Option("icon"),
		Fields:     ctx.OptionSlice("field"),
	}
	spec, err := generator.Normalize(input)
	if err != nil {
		return err
	}
	artifacts, err := generator.RenderResourcePipeline(spec, time.Now().UTC().Format("20060102150405"))
	if err != nil {
		return err
	}
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	runtimeArtifacts, err := generator.RenderRuntimeRegistration(root, spec)
	if err != nil {
		return err
	}
	artifacts = append(artifacts, runtimeArtifacts...)
	if err := generator.WriteAll(root, artifacts); err != nil {
		return err
	}

	ctx.Success("Generated resource skeleton:")
	for _, artifact := range artifacts {
		ctx.Line("- " + strings.ReplaceAll(artifact.Path, "\\", "/"))
	}
	return nil
}
