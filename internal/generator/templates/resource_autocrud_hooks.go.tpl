package resource_hooks

func (r *{{ .ResourceConfig.Package }}.{{ .ResourceConfig.Kind }}) BeforeCreate(m *{{ .ResourceConfig.Package }}.{{ .ResourceConfig.Kind }}Model) {
	// TODO: Add BeforeCreate logic
}

func (r *{{ .ResourceConfig.Package }}.{{ .ResourceConfig.Kind }}) AfterCreate(m *{{ .ResourceConfig.Package }}.{{ .ResourceConfig.Kind }}Model) {
    // TODO: Add AfterCreate logic
}