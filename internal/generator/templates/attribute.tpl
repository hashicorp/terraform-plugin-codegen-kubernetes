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
Default: {{- if eq .AttributeType "StringAttribute"}} "THIS IS A STRING" {{- end}} {{- if eq .AttributeType "BoolAttribute"}} "THIS IS A BOOL" {{- end}} {{- if eq .AttributeType "Int64Attribute"}} "THIS IS An INT64" {{- end}}
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
