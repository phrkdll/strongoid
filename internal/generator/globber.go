package generator

import "path/filepath"

type Globber interface {
	Glob(pattern string) ([]string, error)
}

type RealGlobber struct{}

func (RealGlobber) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
