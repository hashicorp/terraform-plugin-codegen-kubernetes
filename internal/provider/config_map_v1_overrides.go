// THIS FILE HAS BEEN GENERATED FOR YOU TO EDIT
//
// This file allows you to override, intercept, and mutate the behavior
// of the code generated for the resource.
//
// This code was written by a robot on Apr 21, 2023 17:43:05 UTC.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *ConfigMap) overrideSchema(schema schema.Schema) schema.Schema {
	// You can mutate the schema here
	return schema
}

func (r *ConfigMap) afterCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (r *ConfigMap) beforeCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (r *ConfigMap) beforeRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ConfigMap) afterRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ConfigMap) beforeUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ConfigMap) afterUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ConfigMap) beforeDelete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *ConfigMap) afterDelete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
