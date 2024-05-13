// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generator

import (
	"strings"
)

// explicitFields contains a mapping of common snake_case terraform attributes
// and their mapping to Kubernetes manifest fields, and Model struct fields
var explicitMappings = map[string]struct {
	Kubernetes string
	Model      string
}{
	"api_version": {"apiVersion", "APIVersion"},
	"uid":         {"uid", "UID"},
	"id":          {"id", "ID"},
}

// MapTerraformAttributeToKubernetes maps a string containing snake_case into camelCase
func MapTerraformAttributeToKubernetes(terraformAttributeName string) string {
	if v, ok := explicitMappings[terraformAttributeName]; ok {
		return v.Kubernetes
	}

	out := ""
	cap := false
	for _, ch := range terraformAttributeName {
		if ch == '_' {
			cap = true
			continue
		}
		if cap {
			out += strings.ToUpper(string(ch))
			cap = false
		} else {
			out += string(ch)
		}
	}
	return out
}

func MapTerraformAttributeToModel(terraformAttributeName string) string {
	if v, ok := explicitMappings[terraformAttributeName]; ok {
		return v.Model
	}

	// FIXME strings.Title is deprecated
	return strings.Title(MapTerraformAttributeToKubernetes(terraformAttributeName))
}
