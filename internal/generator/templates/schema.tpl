schema.Schema{
    MarkdownDescription: `{{ .Description }}`,
    Blocks: map[string]schema.Block{
        "timeouts": timeouts.BlockAll(ctx),
    },
    {{ .Attributes }}
}