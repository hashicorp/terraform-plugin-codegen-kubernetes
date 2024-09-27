package {{ .ResourceConfig.Package }}

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-codegen-kubernetes/autocrud"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *{{ .ResourceConfig.Kind }}) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var dataModel {{ .ResourceConfig.Kind }}Model

	diag := req.Config.Get(ctx, &dataModel)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	defaultTimeout, err := time.ParseDuration("{{ .ResourceConfig.Generate.Timeouts.Create }}")
	if err != nil {
		resp.Diagnostics.AddError("Error parsing timeout", err.Error())
	return 
	}
	timeout, diag := dataModel.Timeouts.Create(ctx, defaultTimeout)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook.Create -}}
	r.BeforeCreate(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}

	err = autocrud.Create(ctx, r.clientGetter, r.APIVersion, r.Kind, &dataModel)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
		return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook.Create -}}
	r.AfterCreate(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}

	diags := resp.State.Set(ctx, &dataModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *{{ .ResourceConfig.Kind }}) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var dataModel {{ .ResourceConfig.Kind }}Model

	diag := req.State.Get(ctx, &dataModel)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	defaultTimeout, err := time.ParseDuration("{{ .ResourceConfig.Generate.Timeouts.Read }}")
	if err != nil {
		resp.Diagnostics.AddError("Error parsing timeout", err.Error())
		return 
	}
	timeout, diag := dataModel.Timeouts.Read(ctx, defaultTimeout) 
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook.Read -}}
	r.BeforeRead(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}

	var id string
    req.State.GetAttribute(ctx, path.Root("id"), &id)
	err = autocrud.Read(ctx, r.clientGetter, r.Kind, r.APIVersion, id, &dataModel)
	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook.Read -}}
	r.AfterRead(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
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

	defaultTimeout, err := time.ParseDuration("{{ .ResourceConfig.Generate.Timeouts.Update }}")
	if err != nil {
		resp.Diagnostics.AddError("Error parsing timeout", err.Error())
		return 
	}
	timeout, diag := dataModel.Timeouts.Update(ctx, defaultTimeout)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook.Update -}}
	r.BeforeUpdate(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}

	err = autocrud.Update(ctx, r.clientGetter, r.Kind, r.APIVersion, &dataModel)
	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
		return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook.Update -}}
	r.AfterUpdate(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
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

	var dataModel {{ .ResourceConfig.Kind }}Model

	diag := req.State.Get(ctx, &dataModel)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}

	defaultTimeout, err := time.ParseDuration("{{ .ResourceConfig.Generate.Timeouts.Delete }}")
	if err != nil {
		resp.Diagnostics.AddError("Error parsing timeout", err.Error())
		return 
	}
	timeout, diag := dataModel.Timeouts.Delete(ctx, defaultTimeout)
	resp.Diagnostics.Append(diag...)
	if diag.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.BeforeHook.Delete -}}
	r.BeforeDelete(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}

	err = autocrud.Delete(ctx, r.clientGetter, r.Kind, r.APIVersion, req, waitForDeletion)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
		return
	}

	{{ if .ResourceConfig.Generate.CRUDAutoOptions -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook -}}
	{{ if .ResourceConfig.Generate.CRUDAutoOptions.Hooks.AfterHook.Delete -}}
	r.AfterDelete(ctx, req, resp, &dataModel)
	{{ end }}
	{{ end }}
	{{ end }}
	{{ end }}
}

func (r *{{ .ResourceConfig.Kind }}) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var dataModel {{ .ResourceConfig.Kind }}Model

	err := autocrud.Read(ctx, r.clientGetter, r.Kind, r.APIVersion, req.ID, &dataModel)
	if err != nil {
		resp.Diagnostics.AddError("Error importing resource", err.Error())
		return
	}

	// awkward timeouts/types.Object issue https://github.com/hashicorp/terraform-plugin-framework-timeouts/issues/46 & https://github.com/hashicorp/terraform-plugin-framework/issues/716
	dataModel.Timeouts = timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"delete": types.StringType,
			"read":   types.StringType,
			"update": types.StringType,
		}),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataModel)...)
}
