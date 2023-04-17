MarkdownDescription: `{{ .Description }}`,
{{- if .ElementType }}
ElementType: types.{{ .ElementType }},
{{- end }}
{{- if .Required }}
Required: true,
{{- else }}
Optional: true,
{{- end }}