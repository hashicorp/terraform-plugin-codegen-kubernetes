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

{{- if .Sensitive }}
Sensitive: true,
{{- end }}

{{- if .Immutable }}
PlanModifiers: []planmodifier.{{ .PlanModifierType }}{
  {{ .PlanModifierPackage }}.RequiresReplace(), 
},
{{- end }}

{{/* TODO don't share PlanModifierType */}}
{{- if ne .GenAIValidatorType "" }}
Validators: []validator.{{ .PlanModifierType }}{
  {{ .GenAIValidatorType }}{}, 
},
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
