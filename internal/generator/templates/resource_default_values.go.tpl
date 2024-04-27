package {{ .ResourceConfig.Package }}

import (
{{$intExists := "" | and ($stringExists := "") |  and ($boolExists := "")}}
    {{- range .Schema.Attributes }}
{{- if (eq .AttributeType "StringAttribute" | and (eq $stringExists "")) | or (eq .AttributeType "BoolAttribute" | and (eq $boolExists "")) | or (eq .AttributeType "Int64Attribute" | and (eq $intExists "")) | and .DefaultValue }}
    	"github.com/hashicorp/terraform-plugin-framework/resource/schema/{{- if eq .AttributeType "StringAttribute"}}stringdefault{{$stringExists = "true"}}{{end}} {{- if eq .AttributeType "BoolAttribute"}}booldefault{{$boolExists = "true"}}{{end}} {{- if eq .AttributeType "Int64Attribute"}}int64default{{$intExists = "true"}}{{end}}"
{{- end}}
{{- end}}

)

{{- range .Schema.Attributes }}
{{- if eq .AttributeType "StringAttribute" | or (eq .AttributeType "BoolAttribute") | or (eq .AttributeType "Int64Attribute") | and .DefaultValue}}
var {{.Name}}DefaultValue = {{- if eq .AttributeType "StringAttribute"}}stringdefault.StaticString("") // TODO: add default value{{- end}} {{- if eq .AttributeType "BoolAttribute"}}booldefault.StaticBool(false) //TODO change to default value {{- end}} {{- if eq .AttributeType "Int64Attribute"}}int64default.StaticInt64(0) //TODO change to default value{{- end}}
{{- end}}
{{- end}}