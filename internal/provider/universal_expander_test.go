package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestUniversalExpand(t *testing.T) {
	in := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":         "test",
			"namespace":    "kube-system",
			"generateName": "ignore-this",
			"labels": map[string]interface{}{
				"test": "test",
			},
			"annotations": map[string]interface{}{
				"test": "test",
			},
		},
		"immutable": false,
		"data": map[string]interface{}{
			"PGUSER": "hello",
		},
		"binaryData": map[string]interface{}{
			"OBJECT": "aGVsbG8gd29ybGQ=",
		},
	}

	stringMapType := tftypes.Map{ElementType: tftypes.String}
	metadataType := tftypes.List{
		ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"name":        tftypes.String,
				"namespace":   tftypes.String,
				"labels":      stringMapType,
				"annotations": stringMapType,
			},
		},
	}

	ignoredFields := []string{"generateName"}
	actual := UniversalExpand(in, ignoredFields)
	expected := tftypes.NewValue(
		tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"api_version": tftypes.String,
				"kind":        tftypes.String,
				"metadata":    metadataType,
				"data":        stringMapType,
				"binary_data": stringMapType,
				"immutable":   tftypes.Bool,
			},
		},
		map[string]tftypes.Value{
			"api_version": tftypes.NewValue(tftypes.String, "v1"),
			"kind":        tftypes.NewValue(tftypes.String, "ConfigMap"),
			"metadata": tftypes.NewValue(metadataType, []tftypes.Value{
				tftypes.NewValue(metadataType.ElementType, map[string]tftypes.Value{
					"name":      tftypes.NewValue(tftypes.String, "test"),
					"namespace": tftypes.NewValue(tftypes.String, "kube-system"),
					"annotations": tftypes.NewValue(stringMapType, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "test"),
					}),
					"labels": tftypes.NewValue(stringMapType, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "test"),
					}),
				}),
			}),
			"immutable": tftypes.NewValue(tftypes.Bool, false),
			"data": tftypes.NewValue(stringMapType, map[string]tftypes.Value{
				"PGUSER": tftypes.NewValue(tftypes.String, "hello"),
			}),
			"binary_data": tftypes.NewValue(stringMapType, map[string]tftypes.Value{
				"OBJECT": tftypes.NewValue(tftypes.String, "aGVsbG8gd29ybGQ="),
			}),
		},
	)

	assert.Equal(t, expected, actual)
}

func TestSnakify(t *testing.T) {
	cases := []struct {
		camel string
		snake string
	}{
		{
			camel: "apiVersion",
			snake: "api_version",
		},
		{
			camel: "containerPort",
			snake: "container_port",
		},
		{
			camel: "binaryData",
			snake: "binary_data",
		},
		{
			camel: "targetPort",
			snake: "target_port",
		},
		{
			camel: "livenessProbe",
			snake: "liveness_probe",
		},
		{
			camel: "clusterIP",
			snake: "cluster_ip",
		},
		{
			camel: "externalIPs",
			snake: "external_ips",
		},
		{
			camel: "podCIDR",
			snake: "pod_cidr",
		},
	}

	for _, c := range cases {
		actual := Snakify(c.camel)
		assert.Equal(t, c.snake, actual)
	}
}
