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
	// TODO support generic defaults
	DefaultNamespace bool
	Imports          []string
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

// Recursively marks attributes that are metadata.namespace as a quick and dirty
// fix for having namespace support a default value of "default".
// Returns true if it has found any matching attribute in the resource.
// TODO handle lists?
func metadataNamespaceDefault(a AttributesGenerator, parent string) bool {
	for i := range a {
		if a[i].Name == "namespace" && parent == "metadata" {
			a[i].DefaultNamespace = true
			return true
		} else if len(a[i].NestedAttributes) > 0 {
			return metadataNamespaceDefault(a[i].NestedAttributes, a[i].Name)
		}
	}
	return false
}
