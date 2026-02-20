// Copyright IBM Corp. 2024
// SPDX-License-Identifier: MPL-2.0

package generator

import "testing"

func TestMapTerraformAttributeToKubernetes(t *testing.T) {
	testCases := map[string]string{
		"api_version":        "apiVersion",
		"metadata":           "metadata",
		"uid":                "uid",
		"id":                 "id",
		"kind":               "kind",
		"creation_timestamp": "creationTimestamp",
	}

	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			actual := MapTerraformAttributeToKubernetes(input)
			if actual != expected {
				t.Fatalf("expected %q got %q", expected, actual)
			}
		})
	}
}

func TestMapTerraformAttributeModel(t *testing.T) {
	testCases := map[string]string{
		"api_version":        "APIVersion",
		"metadata":           "Metadata",
		"uid":                "UID",
		"id":                 "ID",
		"kind":               "Kind",
		"creation_timestamp": "CreationTimestamp",
	}

	for input, expected := range testCases {
		t.Run(input, func(t *testing.T) {
			actual := MapTerraformAttributeToModel(input)
			if actual != expected {
				t.Fatalf("expected %q got %q", expected, actual)
			}
		})
	}
}
