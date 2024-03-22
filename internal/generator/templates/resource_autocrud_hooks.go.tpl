package {{ .ResourceConfig.Package }}

{{ if .ResourceConfig.Generate.CRUDAutoHooks.BeforeCreate -}}
func (r *{{ .ResourceConfig.Kind }}) BeforeCreate(m *{{ .ResourceConfig.Kind }}Model) {
	// TODO: Add BeforeCreate logic
}
{{ end }}

{{ if .ResourceConfig.Generate.CRUDAutoHooks.AfterCreate -}}
func (r *{{ .ResourceConfig.Kind }}) AfterCreate(m *{{ .ResourceConfig.Kind }}Model) {
    // TODO: Add AfterCreate logic
}
{{ end }}