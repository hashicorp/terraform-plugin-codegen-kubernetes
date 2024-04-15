{{- range $val := .ResourceConfig.DefaultValueAttributes }}
var $val = test
{{- end }}

{{- if eq .AttributeType "StringAttribute"}}stringdefault.StaticString(""), // TODO: add default value{{- end}} {{- if eq .AttributeType "BoolAttribute"}}booldefault.StaticBool(false), //TODO change to default value {{- end}} {{- if eq .AttributeType "Int64Attribute"}}int64default.StaticInt64(0), //TODO change to default value{{- end}}
