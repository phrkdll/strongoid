package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

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
	var logBuf bytes.Buffer

	// Simple testing template
	tmpl := `{{ range .Types }}
	func (id {{ .Name }}) String() string { return fmt.Sprintf("ID: %v", id) }
	{{ end }}`

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
		Generate(
			"testdata",
			[]string{tmpl},
			[]string{"fmt"},
			writer,
			parser,
			glob,
			&logBuf,
		)
	})

	// Check which file was written
	if len(writer.Files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(writer.Files))
	}

	// Check content
	for name, content := range writer.Files {
		if !strings.Contains(string(content), "func (id UserId) String()") {
			t.Errorf("Generated content for %s does not contain expected function", name)
		}
		if !strings.Contains(string(content), "fmt") {
			t.Errorf("Expected import 'fmt' missing in generated file %s", name)
		}
	}

	// Make log output visible
	t.Logf("LOG:\n%s", logBuf.String())
}
