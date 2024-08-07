// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generator

type AttributeGenerator struct {
	Name string

	AttributeType string
	ElementType   string

	PlanModifierType    string
	PlanModifierPackage string

	Required           bool
	Description        string
	Computed           bool
	Sensitive          bool
	Immutable          bool
	GenAIValidatorType string

	NestedAttributes AttributesGenerator
}

func (g AttributeGenerator) String() string {
	return renderTemplate(attributeTemplate, g)
}

type AttributesGenerator []AttributeGenerator

func (g AttributesGenerator) String() string {
	return renderTemplate(attributesTemplate, g)
}
