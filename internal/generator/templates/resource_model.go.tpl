package {{ .ResourceConfig.Package }}

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type {{ .ResourceConfig.Kind }}Model struct {
  Timeouts    timeouts.Value `tfsdk:"timeouts"`
  {{ .ModelFields }}
}
