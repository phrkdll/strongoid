package generator_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/phrkdll/strongoid/internal/generator"
	"github.com/stretchr/testify/assert"
)

type mockFileWriter struct {
	Files map[string][]byte
}

func (m *mockFileWriter) WriteFile(filename string, content []byte) error {
	if m.Files == nil {
		m.Files = map[string][]byte{}
	}
	m.Files[filename] = content
	return nil
}

type mockParser struct {
	Sources map[string]string
}

func (p mockParser) ParseFile(fset *token.FileSet, filename string, _ any) (*ast.File, error) {
	src := p.Sources[filename]
	return parser.ParseFile(fset, filename, src, 0)
}

type mockGlobber struct {
	Files []string
}

func (g mockGlobber) Glob(pattern string) ([]string, error) {
	return g.Files, nil
}

func TestGenerate(t *testing.T) {

	// Simulated file list
	glob := mockGlobber{
		Files: []string{"user.go"},
	}

	writer := &mockFileWriter{}

	parser := mockParser{
		Sources: map[string]string{
			"user.go": `
			package example
			import "strongoid"
			type UserId strongoid.Id[int64]
		`,
		},
	}

	assert.NotPanics(t, func() {
		generator.Generate(
			"testdata",
			[]string{"gorm", "json"},
			writer,
			parser,
			glob,
		)
	})

	// Check which file was written
	if len(writer.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(writer.Files))
	}

	// Check content
	for name, content := range writer.Files {
		if !strings.Contains(string(content), "func (t *UserId) Scan(dbValue any) error") {
			t.Errorf("Generated content for %s does not contain expected function", name)
		}
	}
}
func TestGenerate_SkipsTestAndGenFiles(t *testing.T) {
	glob := mockGlobber{
		Files: []string{"user.go", "user_test.go", "user.gen.go"},
	}
	writer := &mockFileWriter{}
	parser := mockParser{
		Sources: map[string]string{
			"user.go": `
				package example
				import "strongoid"
				type UserId strongoid.Id[int64]
			`,
			"user_test.go": `
				package example
				import "strongoid"
				type ShouldSkip strongoid.Id[int64]
			`,
			"user.gen.go": `
				package example
				import "strongoid"
				type ShouldSkipGen strongoid.Id[int64]
			`,
		},
	}
	generator.Generate(
		"testdata",
		[]string{"gorm", "json"},
		writer,
		parser,
		glob,
	)
	if len(writer.Files) != 1 {
		t.Errorf("Expected only 1 file to be generated, got %d", len(writer.Files))
	}
	for name := range writer.Files {
		if name != "user.gen.go" {
			t.Errorf("Expected generated file to be user.gen.go, got %s", name)
		}
	}
}

func TestGenerate_SkipsInvalidAST(t *testing.T) {
	glob := mockGlobber{
		Files: []string{"broken.go"},
	}
	writer := &mockFileWriter{}
	parser := mockParser{
		Sources: map[string]string{
			"broken.go": `
				package example
				type NotStrongOid int
			`,
		},
	}
	generator.Generate(
		"testdata",
		[]string{"gorm"},
		writer,
		parser,
		glob,
	)
	if len(writer.Files) != 0 {
		t.Errorf("Expected no files to be generated for invalid AST, got %d", len(writer.Files))
	}
}

func TestGenerate_UnsupportedPointerType(t *testing.T) {
	glob := mockGlobber{
		Files: []string{"ptr.go"},
	}
	writer := &mockFileWriter{}
	parser := mockParser{
		Sources: map[string]string{
			"ptr.go": `
				package example
				import "strongoid"
				type PtrId strongoid.Id[*chan int]
			`,
		},
	}
	generator.Generate(
		"testdata",
		[]string{},
		writer,
		parser,
		glob,
	)
	if len(writer.Files) != 0 {
		t.Errorf("Expected no files to be generated for unsupported pointer type, got %d", len(writer.Files))
	}
}

func TestGenerate_SelectorExprBaseType(t *testing.T) {
	glob := mockGlobber{
		Files: []string{"selector.go"},
	}
	writer := &mockFileWriter{}
	parser := mockParser{
		Sources: map[string]string{
			"selector.go": `
				package example
				import (
					"strongoid"
					"foo"
				)
				type FooId strongoid.Id[foo.Bar]
			`,
		},
	}
	generator.Generate(
		"testdata",
		[]string{"json"},
		writer,
		parser,
		glob,
	)
	if len(writer.Files) != 1 {
		t.Errorf("Expected 1 file to be generated for selector expr, got %d", len(writer.Files))
	}
	for _, content := range writer.Files {
		if !strings.Contains(string(content), "FooId") {
			t.Errorf("Generated content does not contain expected type FooId")
		}
	}
}

func TestGenerate_StarSelectorExprBaseType(t *testing.T) {
	glob := mockGlobber{
		Files: []string{"starselector.go"},
	}
	writer := &mockFileWriter{}
	parser := mockParser{
		Sources: map[string]string{
			"starselector.go": `
				package example
				import (
					"strongoid"
					"foo"
				)
				type FooPtrId strongoid.Id[*foo.Bar]
			`,
		},
	}
	generator.Generate(
		"testdata",
		[]string{"json"},
		writer,
		parser,
		glob,
	)
	if len(writer.Files) != 1 {
		t.Errorf("Expected 1 file to be generated for star selector expr, got %d", len(writer.Files))
	}
	for _, content := range writer.Files {
		if !strings.Contains(string(content), "FooPtrId") {
			t.Errorf("Generated content does not contain expected type FooPtrId")
		}
	}
}

func TestGenerate_UsesCorrectTemplates(t *testing.T) {
	glob := mockGlobber{
		Files: []string{"user.go"},
	}
	writer := &mockFileWriter{}
	parser := mockParser{
		Sources: map[string]string{
			"user.go": `
				package example
				import "strongoid"
				type UserId strongoid.Id[int64]
			`,
		},
	}
	generator.Generate(
		"testdata",
		[]string{"gorm", "json"},
		writer,
		parser,
		glob,
	)
	foundGorm := false
	foundJson := false
	for _, content := range writer.Files {
		if strings.Contains(string(content), "Scan(dbValue any) error") {
			foundGorm = true
		}
		if strings.Contains(string(content), "MarshalJSON") {
			foundJson = true
		}
	}
	if !foundGorm {
		t.Errorf("Expected generated file to contain GORM methods")
	}
	if !foundJson {
		t.Errorf("Expected generated file to contain JSON methods")
	}
}
