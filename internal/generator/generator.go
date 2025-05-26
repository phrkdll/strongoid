package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
)

type TypeInfo struct {
	Name     string // i.e. "UserId"
	BaseType string // i.e. "int64", "string"
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
	var foundTypes []TypeInfo

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

					idxExpr, ok := ts.Type.(*ast.IndexExpr)
					if !ok {
						continue
					}

					// Check selector expression: strongoid.Id
					selExpr, ok := idxExpr.X.(*ast.SelectorExpr)
					if !ok || selExpr.Sel.Name != "Id" {
						continue
					}
					pkgIdent, ok := selExpr.X.(*ast.Ident)
					if !ok || pkgIdent.Name != "strongoid" {
						continue
					}

					// Base type: the [T] part
					var baseType string
					switch bt := idxExpr.Index.(type) {
					case *ast.Ident:
						baseType = bt.Name
					default:
						continue
					}

					foundTypes = append(foundTypes, TypeInfo{
						Name:     ts.Name.Name,
						BaseType: baseType,
					})
				}
			}
		}
	}

	if len(foundTypes) == 0 {
		fmt.Println("No strongoid.Id-based types found.")
		return
	}

	f, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	tmpl := template.Must(template.New("json").Parse(methodTemplate))
	err = tmpl.Execute(f, struct {
		Package string
		Types   []TypeInfo
	}{
		Package: pkgName,
		Types:   foundTypes,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Generated '" + outputFileName + "' for IDs:")
	for _, t := range foundTypes {
		fmt.Printf("- %s based on %s\n", t.Name, t.BaseType)
	}
}
