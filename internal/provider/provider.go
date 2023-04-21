package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure KubernetesProvider satisfies various provider interfaces.
var _ provider.Provider = &KubernetesProvider{}

// KubernetesProvider defines the provider implementation.
type KubernetesProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KubernetesProviderModel describes the provider data model.
type KubernetesProviderModel struct {
	ConfigPath string `tfsdk:"config_path"`
}

func (p *KubernetesProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kubernetesx"
	resp.Version = p.version
}

func (p *KubernetesProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config_path": schema.StringAttribute{
				MarkdownDescription: "Path to kubeconfig",
				Required:            true,
			},
		},
	}
}

func (p *KubernetesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KubernetesProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clientGetter := NewKubernetesClientGetter(data.ConfigPath)
	resp.ResourceData = &clientGetter
}

func (p *KubernetesProvider) Resources(ctx context.Context) []func() resource.Resource {
	return resources
}

func (p *KubernetesProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KubernetesProvider{
			version: version,
		}
	}
}
