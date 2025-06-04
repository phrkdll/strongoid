package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Parser interface {
	ParseFile(fset *token.FileSet, filename string, src any) (*ast.File, error)
}

type RealParser struct{}

func (p RealParser) ParseFile(fset *token.FileSet, filename string, src any) (*ast.File, error) {
	return parser.ParseFile(fset, filename, src, 0)
}
