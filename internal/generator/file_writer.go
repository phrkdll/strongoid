package generator

import "os"

type FileWriter interface {
	WriteFile(filename string, content []byte) error
}

type OSFileWriter struct{}

func (OSFileWriter) WriteFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, 0644)
}
