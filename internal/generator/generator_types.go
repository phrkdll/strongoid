package generator

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
