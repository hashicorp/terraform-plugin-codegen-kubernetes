MarkdownDescription: `{{ .Description }}`,
{{- if .ElementType }}
ElementType: types.{{ .ElementType }},
{{- end }}
{{- if .Required }}
Required: true,
{{- else }}
Optional: true,
{{- end }}
{{- if .Computed }}
Computed: true,
{{- end }}
{{- if .DefaultValue}}
Default: {{- if eq .AttributeType "StringAttribute"}}stringdefault.StaticString(""), // TODO: add default value{{- end}} {{- if eq .AttributeType "BoolAttribute"}}booldefault.StaticBool(false), //TODO change to default value {{- end}} {{- if eq .AttributeType "Int64Attribute"}}int64default.StaticInt64(0), //TODO change to default value{{- end}}
{{- end}}
{{- if .Sensitive }}
Sensitive: true,
{{- end }}
{{- if .NestedAttributes }}
  {{- if eq .AttributeType "ListNestedAttribute" }}
  NestedObject: schema.NestedAttributeObject{
    {{ .NestedAttributes }}
  },
  {{- else }}
  {{ .NestedAttributes }}
  {{- end }}
{{- end }}
