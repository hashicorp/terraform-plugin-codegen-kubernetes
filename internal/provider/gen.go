package provider

//go:generate go run ../../tools/codegen/openapi_v3.go -pkg provider -json ../../tools/codegen/data/k8s-1.24.3-core-v1.json -ref io.k8s.api.core.v1.ConfigMap -o config_map_resource.go
