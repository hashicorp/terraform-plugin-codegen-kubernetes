// Code generated by tools/gen/schema_openapi_v3.go; DO NOT EDIT 
//
// This file contains the autogenerated implementation for the "{{ .ResourceName }}" resource. 
// You can override the behaviour of this resource using the corresponding _overrides.go file. 
//
// This code was written by a robot on {{ .GeneratedTimestamp.Format "Jan 02, 2006 15:04:05 UTC" }}.

package {{ .Package }}

import (
	"fmt"
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &{{ .Kind }}{}
var _ resource.ResourceWithImportState = &{{ .Kind }}{}

type {{ .Kind }} struct{
	Kind string
	APIVersion string

	IgnoredFields []string

	clientGetter *KubernetesClientGetter
}

func New{{ .Kind }}() resource.Resource {
	return &{{ .Kind }}{
		Kind: "{{ .Kind }}",
		APIVersion: "{{ .APIVersion }}",
		IgnoredFields: []string{
			{{- range $v := .IgnoredFields }}
			"{{ $v }}",
			{{- end }}
		},
	}
}

func (r *{{ .Kind }}) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{ .ResourceName }}"
}

func (r *{{ .Kind }}) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `{{ .Root.Description }}`,
		Attributes:   map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			{{- range $val := .Root.Attributes }}
			"{{- $val.Name }}": schema.{{ $val.AttributeType }}{ 
				{{ $val }} 
			},
			{{- end }}
		},
		Blocks:  map[string]schema.Block{
			{{- range $val := .Root.Blocks }}
			"{{- $val.Name }}": schema.ListNestedBlock{ 
				{{ $val }} 
			},
			{{- end }}
		},
	}
}

func (r *{{ .Kind }}) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientGetter, ok := req.ProviderData.(*KubernetesClientGetter)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *KubernetesClientGetter, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.clientGetter = clientGetter
}


func (r *{{ .Kind }}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state, err := UniversalCreate(ctx, r.clientGetter, r.Kind, r.APIVersion, r.IgnoredFields, req)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
	}

	resp.State.Raw = state
}

func (r *{{ .Kind }}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state, err := UniversalRead(ctx, r.clientGetter, r.Kind, r.APIVersion, r.IgnoredFields, req)
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
	}

	resp.State.Raw = state
}

func (r *{{ .Kind }}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	state, err := UniversalUpdate(ctx, r.clientGetter, r.Kind, r.APIVersion, r.IgnoredFields, req)
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
	}

	resp.State.Raw = state
}

func (r *{{ .Kind }}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	err := UniversalDelete(ctx, r.clientGetter, r.Kind, r.APIVersion, req)
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
	}
}

func (r *{{ .Kind }}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}