package main

import (
	"flag"
	"fmt"
	"os"
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

	generator.Generate(".", modules, generator.OSFileWriter{}, generator.RealParser{}, generator.RealGlobber{})
}
