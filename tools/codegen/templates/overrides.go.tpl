// THIS FILE HAS BEEN GENERATED FOR YOU TO EDIT
//
// This file allows you to override, intercept, and mutate the behavior 
// of the code generated for the resource.
//
// This code was written by a robot on {{ .GeneratedTimestamp.Format "Jan 02, 2006 15:04:05 UTC" }}.

package {{ .Package }}

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func (r *{{ .Kind }}) overrideSchema(schema schema.Schema) schema.Schema {
    // You can mutate the schema here 
    return schema
}

func (r *{{ .Kind }}) afterCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (r *{{ .Kind }}) beforeCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

func (r *{{ .Kind }}) beforeRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *{{ .Kind }}) afterRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *{{ .Kind }}) beforeUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *{{ .Kind }}) afterUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *{{ .Kind }}) beforeDelete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *{{ .Kind }}) afterDelete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}