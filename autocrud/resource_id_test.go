// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package autocrud

import "testing"

func TestCreateID(t *testing.T) {
	cases := map[string]map[string]any{
		"kube-system/example": {
			"metadata": map[string]any{
				"name":      "example",
				"namespace": "kube-system",
			},
		},
		"default/example-resource": {
			"metadata": map[string]any{
				"name":      "example-resource",
				"namespace": "default",
			},
		},
		"example": {
			"metadata": map[string]any{
				"name": "example",
			},
		},
	}

	for id, manifest := range cases {
		t.Run(id, func(t *testing.T) {
			if actualID := createID(manifest); actualID != id {
				t.Fatalf("expected %q got %q", id, actualID)
			}
		})
	}
}

func TestParseID(t *testing.T) {
	cases := map[string]struct {
		namespace string
		name      string
	}{
		"kube-system/example": {
			namespace: "kube-system",
			name:      "example",
		},
		"default/example-resource": {
			namespace: "default",
			name:      "example-resource",
		},
		"example": {
			namespace: "",
			name:      "example",
		},
	}

	for id, expected := range cases {
		t.Run(id, func(t *testing.T) {
			namespace, name := parseID(id)
			if name != expected.name {
				t.Fatalf("expected name %q got %q", expected.name, name)
			}
			if namespace != expected.namespace {
				t.Fatalf("expected namespace %q got %q", expected.namespace, namespace)
			}
		})
	}
}
