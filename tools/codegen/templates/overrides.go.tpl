// THIS FILE IS FOR YOU TO EDIT 
//
// This file allows you to override, intercept, and mutate the behavior 
// of the code generated for the resource.

package {{ .Package }}

func (r *{{ .ResourceName }}) SchemaOverride(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, f func() {}) {
	f()
}
