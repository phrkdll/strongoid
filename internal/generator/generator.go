package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/phrkdll/must/pkg/must"
	tmpls "github.com/phrkdll/strongoid/internal/generator/templates"
)

func Generate(dir string,
	modules []string,
	writer FileWriter,
	parser Parser,
	globber Globber) {

	imports := []string{"github.com/phrkdll/strongoid/pkg/strongoid"}
	templates := []string{tmpls.BaseTemplate}

	if slices.Contains(modules, "gorm") {
		imports = append(imports, "database/sql/driver")
		templates = append(templates, tmpls.GormTemplate)
	}

	if slices.Contains(modules, "json") {
		templates = append(templates, tmpls.JsonTemplate)
	}

	fset := token.NewFileSet()
	files := must.Return(globber.Glob(filepath.Join(dir, "*.go"))).ElsePanic()

	fileMap := map[string]*FileTypes{}

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") || strings.HasSuffix(file, ".gen.go") {
			continue
		}

		node, err := parser.ParseFile(fset, file, nil)
		if err != nil {
			fmt.Printf("Could not parse %s: %v\n", file, err)
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
		var buf bytes.Buffer

		for _, mt := range templates {
			tmpl := template.Must(template.New("strongoid").Parse(mt))

			err := tmpl.Execute(&buf, fileData)
			if err != nil {
				fmt.Printf("Template error on %s: %v\n", genFile, err)
				continue
			}
		}

		err := writer.WriteFile(genFile, buf.Bytes())
		if err != nil {
			fmt.Printf("Failed to write %s: %v\n", genFile, err)
			continue
		}

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
