package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ketryx/pkg/ketryx"
)

var (
	// Ensure provider defined types fully satisfy interfaces
	_ provider.Provider = &ketryxProvider{}
)

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// ketryxProviderModel
//
////////////////////////////////////////////////////////////////////////////////////////////////////

type ketryxProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// ketryxProvider
//
////////////////////////////////////////////////////////////////////////////////////////////////////

type ketryxProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ketryxProvider{
			version: version,
		}
	}
}

// Provider general configuration and naming
func (k *ketryxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ketryx"
	resp.Version = k.version
}

// Provider-level schema for configuration
func (k *ketryxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Prepares a Ketryx API client
func (k *ketryxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var (
		api_key string
		client  *ketryx.Client
		config  ketryxProviderModel
		err     error
	)

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Ketryx API key",
			"Unknown Ketryx API key (longer description)",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	api_key = os.Getenv("KETRYX_API_KEY")
	if !config.APIKey.IsNull() {
		api_key = config.APIKey.ValueString()
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Empty Ketryx API key",
			"Empty Ketryx API key (longer description)",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client, err = ketryx.NewClient(api_key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Ketryx API client error",
			"Ketryx API client error: "+err.Error(),
		)
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// Defines `data` blocks implemented by the provider
func (k *ketryxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewProjectsDataSource,
	}
}

// Defines `resource` blocks implemented by the provider
func (k *ketryxProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
