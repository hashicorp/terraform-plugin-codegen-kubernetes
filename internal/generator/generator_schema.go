// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generator

import (
	"path"
)

type SchemaGenerator struct {
	Name        string
	Description string
	Attributes  AttributesGenerator
	Imports     []string
}

func (g SchemaGenerator) String() string {
	return renderTemplate(schemaTemplate, g)
}

const schemaImportPath = "github.com/hashicorp/terraform-plugin-framework/resource/schema"

// FIXME this should probably be a reciever function
func getPlanModifierImports(a AttributesGenerator) []string {
	imports := []string{}
	for _, aa := range a {
		if aa.Immutable {
			imports = append(imports, path.Join(schemaImportPath, aa.PlanModifierPackage))
		}
		imports = append(imports, getPlanModifierImports(aa.NestedAttributes)...)
	}

	// FIXME: this will need to dedupe
	return imports
}
