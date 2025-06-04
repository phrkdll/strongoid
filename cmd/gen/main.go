package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/phrkdll/strongoid/internal/generator"
	tmpls "github.com/phrkdll/strongoid/internal/generator/templates"
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
	templates := []string{tmpls.BaseTemplate}

	if slices.Contains(modules, "gorm") {
		imports = append(imports, "database/sql/driver")
		templates = append(templates, tmpls.GormTemplate)
	}

	if slices.Contains(modules, "json") {
		templates = append(templates, tmpls.JsonTemplate)
	}

	generator.Generate(".", templates, imports, generator.OSFileWriter{}, generator.RealParser{}, generator.RealGlobber{})
}
