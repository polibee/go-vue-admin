package main

import (
	"flag"
	"fmt"
	"os"

	"goravel/app/core/generator"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "module":
		if err := runModule(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "admin-gen module: %v\n", err)
			os.Exit(1)
		}
	case "resource":
		fmt.Fprintln(os.Stderr, "admin-gen resource: not implemented yet; use the Resource Manifest design in docs/RESOURCE_GENERATION.md")
		os.Exit(2)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "admin-gen: unknown command %q\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func runModule(args []string) error {
	flags := flag.NewFlagSet("module", flag.ContinueOnError)
	root := flags.String("root", ".", "repository root")
	moduleRoot := flags.String("module-root", "", "module output directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("usage: admin-gen module [--root <dir>] [--module-root <dir>] <name>")
	}
	return generator.GenerateModule(generator.ModuleOptions{
		RootDir:    *root,
		ModuleRoot: *moduleRoot,
		Name:       flags.Arg(0),
	})
}

func usage() {
	fmt.Println("admin-gen - Go Vue Admin code generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  admin-gen module <name>")
	fmt.Println("  admin-gen resource --module <module> --table <table>")
}
