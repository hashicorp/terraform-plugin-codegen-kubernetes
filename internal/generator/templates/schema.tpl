schema.Schema{
    MarkdownDescription: `{{ .Description }}`,
    {{- if not .WithoutTimeouts }}
    Blocks: map[string]schema.Block{
        "timeouts": timeouts.BlockAll(),
    },
    {{- end }}
    {{ .Attributes }}
}