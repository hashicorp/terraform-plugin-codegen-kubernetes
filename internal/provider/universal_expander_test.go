package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniversalExpand(t *testing.T) {
	in := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      "test",
			"namespace": "kube-system",
		},
		"immutable": false,
		"binaryData": map[string]interface{}{
			"OBJECT": "aGVsbG8gd29ybGQ=",
		},
		"data": map[string]interface{}{
			"PGUSER": "hello",
		},
	}

	out := ExpandValue(in)

	t.Logf("Output: %#v", out)
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
