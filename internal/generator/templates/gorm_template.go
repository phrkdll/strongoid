package templates

const GormTemplate = `{{ range .Types }}
func (t *{{ .Name }}) Scan(dbValue any) error {
	return (*strongoid.Id[{{ .BaseType }}])(t).Scan(dbValue)
}

func (t {{ .Name }}) Value() (driver.Value, error) {
	return strongoid.Id[{{ .BaseType }}](t).Value()
}
{{ end }}`
