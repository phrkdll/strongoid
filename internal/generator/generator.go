package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TypeInfo struct {
	Name     string
	BaseType string
}

type FileTypes struct {
	File    string
	Package string
	Imports []string
	Types   []TypeInfo
}

func Generate(methodTemplates, imports []string) {
	dir := "."

	fset := token.NewFileSet()
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		panic(err)
	}

	fileMap := map[string]*FileTypes{}

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") || strings.HasSuffix(file, ".gen.go") {
			continue
		}

		node, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse %s: %v\n", file, err)
			continue
		}

		var collected []TypeInfo
		for _, decl := range node.Decls {
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

				selExpr, ok := idxExpr.X.(*ast.SelectorExpr)
				if !ok || selExpr.Sel.Name != "Id" {
					continue
				}
				pkgIdent, ok := selExpr.X.(*ast.Ident)
				if !ok || pkgIdent.Name != "strongoid" {
					continue
				}

				var baseType string
				switch bt := idxExpr.Index.(type) {
				case *ast.Ident:
					baseType = bt.Name
				case *ast.SelectorExpr:
					pkg := bt.X.(*ast.Ident).Name
					sel := bt.Sel.Name
					baseType = pkg + "." + sel
				case *ast.StarExpr:
					switch inner := bt.X.(type) {
					case *ast.Ident:
						baseType = "*" + inner.Name
					case *ast.SelectorExpr:
						pkg := inner.X.(*ast.Ident).Name
						sel := inner.Sel.Name
						baseType = "*" + pkg + "." + sel
					default:
						fmt.Printf("Unsupported pointer type in %s\n", ts.Name.Name)
						continue
					}
				default:
					fmt.Printf("Unsupported type in %s\n", ts.Name.Name)
					continue
				}

				collected = append(collected, TypeInfo{
					Name:     ts.Name.Name,
					BaseType: baseType,
				})
			}
		}

		if len(collected) > 0 {
			fileMap[file] = &FileTypes{
				File:    file,
				Package: node.Name.Name,
				Imports: imports,
				Types:   collected,
			}
		}
	}

	// Generate for each source file
	for src, fileData := range fileMap {
		genFile := strings.TrimSuffix(src, ".go") + ".gen.go"

		f, err := os.Create(genFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", genFile, err)
			continue
		}

		for _, mt := range methodTemplates {

			tmpl := template.Must(template.New("strongoid").Parse(mt))

			err = tmpl.Execute(f, fileData)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Template error on %s: %v\n", genFile, err)
				_ = f.Close()
				continue
			}
		}

		_ = f.Close()
		fmt.Printf("Generated '%s' for %s in '%s'\n", genFile, joinTypeNames(fileData.Types), src)
	}
}

func joinTypeNames(types []TypeInfo) string {
	var names []string
	for _, t := range types {
		names = append(names, t.Name)
	}
	return strings.Join(names, ", ")
}
