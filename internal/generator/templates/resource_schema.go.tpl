package {{ .ResourceConfig.Package }}

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
  {{- range $val := .Schema.Imports }}
  "{{ $val }}"
  {{- end }}
)

func (r *{{ .ResourceConfig.Kind }}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = {{ .Schema }}
}
