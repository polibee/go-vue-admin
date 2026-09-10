package main

import (
	"context"
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
		if err := runResource(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "admin-gen resource: %v\n", err)
			os.Exit(1)
		}
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

func runResource(args []string) error {
	flags := flag.NewFlagSet("resource", flag.ContinueOnError)
	root := flags.String("root", ".", "repository root")
	module := flags.String("module", "", "target module")
	table := flags.String("table", "", "MySQL table")
	name := flags.String("name", "", "resource id; defaults to table")
	label := flags.String("label", "", "resource label; defaults to inferred table label")
	user := flags.String("user", envOr("DB_USERNAME", "root"), "MySQL username")
	password := flags.String("password", envOr("DB_PASSWORD", ""), "MySQL password")
	address := flags.String("address", envOr("DB_HOST", "127.0.0.1")+":"+envOr("DB_PORT", "3306"), "MySQL host:port")
	database := flags.String("database", envOr("DB_DATABASE", ""), "MySQL database")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *module == "" || *table == "" {
		return fmt.Errorf("usage: admin-gen resource --module <module> --table <table> [--name <id>] [--label <label>]")
	}
	if *database == "" {
		return fmt.Errorf("MySQL database is required; set DB_DATABASE or --database")
	}
	resourceID := *name
	if resourceID == "" {
		resourceID = *table
	}
	introspector := generator.NewMySQLIntrospector(context.Background(), *user, *password, *address, *database)
	defer introspector.Close()
	schema, err := introspector.Inspect(context.Background(), *table)
	if err != nil {
		return err
	}
	manifest := generator.ManifestFromTableSchema(schema)
	manifest.ID = resourceID
	if *label != "" {
		manifest.Label = *label
	}
	return generator.GenerateResource(generator.ResourceOptions{
		RootDir:  *root,
		Module:   *module,
		Manifest: manifest,
	})
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func usage() {
	fmt.Println("admin-gen - Go Vue Admin code generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  admin-gen module <name>")
	fmt.Println("  admin-gen resource --module <module> --table <table>")
}
