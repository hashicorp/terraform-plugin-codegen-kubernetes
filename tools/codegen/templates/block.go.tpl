MarkdownDescription: `{{ .Description }}`,
NestedObject: schema.NestedBlockObject{
	Attributes:   map[string]schema.Attribute{
		{{- range $val := .Attributes }}
		"{{- $val.Name }}": schema.{{ $val.AttributeType }}{ 
			{{ $val }} 
		},
		{{- end }}
	},
	Blocks:  map[string]schema.Block{
		{{- range $val := .Blocks }}
		"{{- $val.Name }}": schema.ListNestedBlock{ 
			{{ $val }} 
		},
		{{- end }}
	},
},