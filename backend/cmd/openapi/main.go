package main

import (
	"flag"
	"log"
	"path/filepath"

	"goravel/app/core/openapi"
)

func main() {
	output := flag.String("output", "../../contracts/openapi/openapi.json", "OpenAPI document output path")
	schemaDir := flag.String("schema-dir", "../../contracts/schemas", "schema output directory")
	flag.Parse()

	if err := openapi.WriteJSON(*output, openapi.BuildDocument()); err != nil {
		log.Fatal(err)
	}
	for name, schema := range openapi.SchemaDocuments() {
		if err := openapi.WriteJSON(filepath.Join(*schemaDir, name+".json"), schema); err != nil {
			log.Fatal(err)
		}
	}
}
