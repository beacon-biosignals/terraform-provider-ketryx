package provider

import (
	"context"
	"fmt"
	_ "net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	_ "github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-ketryx/pkg/ketryx"
)

// Ensure provider defined types fully satisfy interfaces
var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSource              = &projectsDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectsDataSource{}
)

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// projectDataSource
//
////////////////////////////////////////////////////////////////////////////////////////////////////

type projectDataSource struct {
	client *ketryx.Client
}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

func (p *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ketryx.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ketryx.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	p.client = client
}

func (d *projectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// projectsDataSource
//
////////////////////////////////////////////////////////////////////////////////////////////////////

type projectsDataSource struct {
	client *ketryx.Client
}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

func (p *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ketryx.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ketryx.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	p.client = client
}

func (d *projectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}

func (d *projectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (p *projectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state projectsDataSourceModel

	projects, err := p.client.ProjectsGet(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Ketryx Project API",
			err.Error(),
		)
	}

	for _, project := range projects.Projects {
		projectState := projectModel{
			Id:   types.StringValue(project.Id),
			Name: types.StringValue(project.Name),
		}

		state.Projects = append(state.Projects, projectState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// Models
//
////////////////////////////////////////////////////////////////////////////////////////////////////

type projectModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type projectsDataSourceModel struct {
	Projects []projectModel `tfsdk:"projects"`
}
