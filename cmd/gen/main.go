package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/phrkdll/strongoid/internal/generator"
)

var (
	modulesFlag = flag.String("modules", "", "Comma-separated list of modules to generate methods for (e.g. json, gorm)")
)

func main() {
	flag.Parse()

	if *modulesFlag == "" {
		fmt.Fprintln(os.Stderr, "missing --modules flag")
		os.Exit(1)
	}

	modules := strings.Split(*modulesFlag, ",")

	imports := []string{"github.com/phrkdll/strongoid/pkg/strongoid"}
	templates := []string{baseTemplate}

	if slices.Contains(modules, "gorm") {
		imports = append(imports, "database/sql/driver")
		templates = append(templates, gormTemplate)
	}

	if slices.Contains(modules, "json") {
		templates = append(templates, jsonTemplate)
	}

	generator.Generate(templates, imports)
}
