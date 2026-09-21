package console

import (
	"fmt"
	"os"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"goravel/app/generator"
)

type ModuleCheckCommand struct{}

func (ModuleCheckCommand) Signature() string { return "admin:check-module" }

func (ModuleCheckCommand) Description() string {
	return "Check a generated module skeleton without changing the project"
}

func (ModuleCheckCommand) Extend() command.Extend {
	return command.Extend{
		Category:  "admin",
		ArgsUsage: "<name>",
		Arguments: []command.Argument{&command.ArgumentString{Name: "name", Usage: "module name", Required: true}},
	}
}

func (ModuleCheckCommand) Handle(ctx console.Context) error {
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	report, err := generator.CheckModule(root, ctx.ArgumentString("name"))
	if err != nil {
		return err
	}
	if !report.Complete {
		ctx.Warning("Module skeleton is incomplete:")
		for _, missing := range report.Missing {
			ctx.Line("- missing " + strings.ReplaceAll(missing, "\\", "/"))
		}
		return fmt.Errorf("module %q is missing %d generated files", report.Name, len(report.Missing))
	}
	ctx.Success(fmt.Sprintf("Module %q is complete; no files were changed.", report.Name))
	return nil
}
