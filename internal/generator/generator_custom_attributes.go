package generator

type CustomAttributesGenerator struct {
	WaitForRollout               bool
	WaitForDefaultServiceAccount bool
	WaitForLoadBalancer          bool
}

func (g CustomAttributesGenerator) String() string {
	return renderTemplate(customAttributeTemplate, g)
}
