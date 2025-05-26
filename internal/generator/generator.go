package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"strings"
)

type templateData struct {
	Package string
	Types   []string
}

func Generate(methodTemplate, outputFileName string) {
	dir := "." // current directory
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		fmt.Println("parse error:", err)
		os.Exit(1)
	}

	var pkgName string
	var foundTypes []string

	for name, pkg := range pkgs {
		pkgName = name // use directory name as package name
		for _, f := range pkg.Files {
			for _, decl := range f.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					// look for underlying type expressions
					ident, ok := ts.Type.(*ast.IndexExpr)
					if !ok {
						continue
					}
					// check if base type is strongoid.Id
					sel, ok := ident.X.(*ast.SelectorExpr)
					if !ok {
						continue
					}
					if sel.Sel.Name != "Id" {
						continue
					}
					pkgIdent, ok := sel.X.(*ast.Ident)
					if !ok || pkgIdent.Name != "strongoid" {
						continue
					}
					// matched: type MyType strongoid.Id[...]
					foundTypes = append(foundTypes, ts.Name.Name)
				}
			}
		}
	}

	if len(foundTypes) == 0 {
		fmt.Println("No strongoid.Id-based types found.")
		return
	}

	// write generated file
	f, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	tmpl := template.Must(template.New("json").Parse(methodTemplate))
	err = tmpl.Execute(f, templateData{
		Package: pkgName,
		Types:   foundTypes,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated '"+outputFileName+"' for IDs:", strings.Join(foundTypes, ", "))
}
