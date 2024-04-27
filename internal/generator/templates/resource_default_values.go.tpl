package {{ .ResourceConfig.Package }}

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

{{- range .Schema.Attributes }}
{{- if eq .AttributeType "StringAttribute" | or (eq .AttributeType "BoolAttribute") | or (eq .AttributeType "Int64Attribute")}}
var {{.Name}}DefaultValue = {{- if eq .AttributeType "StringAttribute"}}stringdefault.StaticString("") // TODO: add default value{{- end}} {{- if eq .AttributeType "BoolAttribute"}}booldefault.StaticBool(false) //TODO change to default value {{- end}} {{- if eq .AttributeType "Int64Attribute"}}int64default.StaticInt64(0) //TODO change to default value{{- end}}
{{- end}}
{{- end}}
