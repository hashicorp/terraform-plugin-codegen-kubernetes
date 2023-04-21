package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/stretchr/testify/assert"
)

func TestFlattenValue(t *testing.T) {
	metadataType := tftypes.List{
		ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":      tftypes.String,
				"namespace": tftypes.String,
			},
		},
	}
	dataType := tftypes.Map{ElementType: tftypes.String}
	in := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"api_version": tftypes.String,
			"kind":        tftypes.String,
			"metadata":    metadataType,
			"data":        dataType,
		},
	}, map[string]tftypes.Value{
		"api_version": tftypes.NewValue(tftypes.String, "v1"),
		"kind":        tftypes.NewValue(tftypes.String, "ConfigMap"),
		"metadata": tftypes.NewValue(metadataType, []tftypes.Value{
			tftypes.NewValue(metadataType.ElementType, map[string]tftypes.Value{
				"name":      tftypes.NewValue(tftypes.String, "example"),
				"namespace": tftypes.NewValue(tftypes.String, "example"),
			}),
		}),
		"data": tftypes.NewValue(dataType, map[string]tftypes.Value{
			"PGUSER": tftypes.NewValue(tftypes.String, "admin"),
			"PGHOST": tftypes.NewValue(tftypes.String, "localhost"),
		}),
	})

	expected := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      "example",
			"namespace": "example",
		},
		"data": map[string]interface{}{
			"PGUSER": "admin",
			"PGHOST": "localhost",
		},
	}
	actual := FlattenValue(in)
	assert.EqualValues(t, expected, actual)
}

func TestCamelize(t *testing.T) {
	cases := []struct {
		snake string
		camel string
	}{
		{
			snake: "api_version",
			camel: "apiVersion",
		},
		{
			snake: "container_port",
			camel: "containerPort",
		},
		{
			snake: "binary_data",
			camel: "binaryData",
		},
		{
			snake: "target_port",
			camel: "targetPort",
		},
		{
			snake: "liveness_probe",
			camel: "livenessProbe",
		},
	}

	for _, c := range cases {
		actual := Camelize(c.snake)
		assert.Equal(t, c.camel, actual)
	}
}
