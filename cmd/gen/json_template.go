package main

const jsonTemplate = `{{ range .Types }}
func (t {{ .Name }}) MarshalJSON() ([]byte, error) {
	return strongoid.Id[{{ .BaseType }}](t).MarshalJSON()
}

func (t *{{ .Name }}) UnmarshalJSON(data []byte) error {
	return (*strongoid.Id[{{ .BaseType }}])(t).UnmarshalJSON(data)
}
{{ end }}`
