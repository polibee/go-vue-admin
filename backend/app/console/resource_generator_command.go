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
			&command.StringFlag{Name: "page-mode", Usage: "page mode: generic or custom"},
			&command.StringFlag{Name: "permission", Usage: "view permission"},
			&command.StringFlag{Name: "icon", Usage: "menu icon"},
			&command.StringSliceFlag{Name: "action", Usage: "action name:label:kind:permission:batch:payload"},
			&command.StringSliceFlag{Name: "action-field", Usage: "action payload field action:name:label:type:required[:value=Label|value=Label]"},
			&command.StringFlag{Name: "scope", Usage: "data scope: all or own"},
			&command.StringFlag{Name: "owner-field", Usage: "integer field used by own data scope"},
			&command.StringSliceFlag{Name: "field", Usage: "field definition name:type[:required[:value=Label|value=Label]]", Required: true},
			&command.StringSliceFlag{Name: "relation", Usage: "relation name:kind:resource:field:foreign_field:label_field[:selectable]"},
			&command.StringSliceFlag{Name: "form-group", Usage: "form group name:label:columns:field|field"},
			&command.StringSliceFlag{Name: "detail-section", Usage: "detail section name:label:field|field"},
		},
	}
}

func (ResourceGeneratorCommand) Handle(ctx console.Context) error {
	input := generator.Input{
		Name:             ctx.ArgumentString("name"),
		Label:            ctx.Option("label"),
		Route:            ctx.Option("route"),
		PageMode:         ctx.Option("page-mode"),
		Permission:       ctx.Option("permission"),
		Icon:             ctx.Option("icon"),
		ActionSpecs:      ctx.OptionSlice("action"),
		ActionFieldSpecs: ctx.OptionSlice("action-field"),
		Scope:            ctx.Option("scope"),
		OwnerField:       ctx.Option("owner-field"),
		Fields:           ctx.OptionSlice("field"),
		Relations:        ctx.OptionSlice("relation"),
		FormGroups:       ctx.OptionSlice("form-group"),
		Details:          ctx.OptionSlice("detail-section"),
	}
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	artifacts, err := generator.GenerateResource(root, input, time.Now().UTC().Format("20060102150405"))
	if err != nil {
		return err
	}

	ctx.Success("Generated resource skeleton:")
	for _, artifact := range artifacts {
		ctx.Line("- " + strings.ReplaceAll(artifact.Path, "\\", "/"))
	}
	return nil
}
