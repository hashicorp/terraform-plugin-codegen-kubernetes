package {{ .ResourceConfig.Package }}


import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-codegen-kubernetes/autocrud"
)

func (r *{{ .ResourceConfig.Kind }}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var dataModel {{ .ResourceConfig.Kind }}Model

	diag := req.Config.Get(ctx, &dataModel)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeCreate -}}
	r.BeforeCreate(&dataModel)
	{{ end }}

	err := autocrud.Create(ctx, r.clientGetter, r.APIVersion, r.Kind, &dataModel)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
    return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterCreate -}}
	r.AfterCreate(&dataModel)
	{{ end }}
	diags := resp.State.Set(ctx, &dataModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *{{ .ResourceConfig.Kind }}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var dataModel {{ .ResourceConfig.Kind }}Model

	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeRead -}}
	r.BeforeRead(&dataModel)
	{{ end }}

	err := autocrud.Read(ctx, r.clientGetter, r.Kind, r.APIVersion, req, &dataModel)
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
    return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterRead -}}
	r.AfterRead(&dataModel)
	{{ end }}

	diags := resp.State.Set(ctx, &dataModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *{{ .ResourceConfig.Kind }}) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var dataModel {{ .ResourceConfig.Kind }}Model

	diag := req.Config.Get(ctx, &dataModel)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeUpdate -}}
	r.BeforeUpdate(&dataModel)
	{{ end }}

	err := autocrud.Update(ctx, r.clientGetter, r.Kind, r.APIVersion, &dataModel)
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
    return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterUpdate -}}
	r.AfterUpdate(&dataModel)
	{{ end }}

	diags := resp.State.Set(ctx, &dataModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *{{ .ResourceConfig.Kind }}) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
  {{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
  waitForDeletion := {{ .ResourceConfig.Generate.CRUDAutoOptions.WaitForDeletion }}
  {{- else -}}
  waitForDeletion := false
  {{- end }}
	err := autocrud.Delete(ctx, r.clientGetter, r.Kind, r.APIVersion, req, waitForDeletion)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
    return
	}
}

func (r *{{ .ResourceConfig.Kind }}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
