package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-ketryx/pkg/ketryx"
)

// Ensure defined types fully satisfy interfaces
var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

////////////////////////////////////////////////////////////////////////////////////////////////////
//
// projectResource
//
////////////////////////////////////////////////////////////////////////////////////////////////////

// projectResource is the resource implementation.
type projectResource struct {
	client *ketryx.Client
}

// NewOrderResource is a helper function to simplify the provider implementation.
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

func (p *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	var (
		client *ketryx.Client
		ok     bool
	)

	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok = req.ProviderData.(*ketryx.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Soure Configure Type",
			fmt.Sprintf("Expected *ketryx.Client, got: %T. Please report this issue to the provider developer", req.ProviderData),
		)
		return
	}

	p.client = client
}

func (p *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		diags    diag.Diagnostics
		err      error
		plan     projectResourceModel
		response ketryx.KetryxAPIProjectsPostResponse
		request  ketryx.KetryxAPIProjectsPostRequest
	)

	// Retrieve the plan
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request
	request.Name = plan.Name.ValueString()
	for _, repository := range plan.Repository {
		request.Repositories = append(request.Repositories, ketryx.KetryxAPIRepository{
			AuthToken:  repository.AuthToken.ValueString(),
			AuthUser:   repository.AuthUser.ValueString(),
			MainRef:    repository.MainRef.ValueString(),
			ReleaseRef: repository.ReleaseRef.ValueString(),
			Url:        repository.Url.ValueString(),
		})
	}

	// Create the project
	response, err = p.client.ProjectsPost(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating project",
			"Could not create project, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state
	plan.Id = types.StringValue(response.Id)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		diags diag.Diagnostics
		err   error
		state projectResourceModel
	)

	// Retrieve current state
	diags = resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the project
	_, err = p.client.ProjectDelete(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting project",
			"Could not delete project, unexpected error: "+err.Error(),
		)
		return
	}
}

func (p *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (p *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (p *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		diags    diag.Diagnostics
		err      error
		response ketryx.KetryxAPIProjectGetResponse
		state    projectResourceModel
	)

	// Retrieve current state
	diags = resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the project
	response, err = p.client.ProjectGet(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving project",
			"Could not get project, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(response.Name)
	for _, repository := range response.Repositories {
		state.Repository = append(state.Repository, projectRepositoryResourceModel{
			AuthToken:  types.StringValue(repository.AuthToken),
			AuthUser:   types.StringValue(repository.AuthUser),
			MainRef:    types.StringValue(repository.MainRef),
			ReleaseRef: types.StringValue(repository.ReleaseRef),
			Url:        types.StringValue(repository.Url),
		})
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"repository": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"auth_token": schema.StringAttribute{
						Required: true,
					},
					"auth_user": schema.Float64Attribute{
						Required: true,
					},
					"main_ref": schema.StringAttribute{
						Required: true,
					},
					"release_ref": schema.StringAttribute{
						Required: true,
					},
					"url": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (p *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		diags   diag.Diagnostics
		err     error
		plan    projectResourceModel
		request ketryx.KetryxAPIProjectPostRequest
		state   projectResourceModel
	)

	// Retrieve current state
	diags = resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the plan
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request
	// request.Settings = plan.Settings.ValueString() // TODO Figure out settings field
	for _, repository := range plan.Repository {
		request.Repositories = append(request.Repositories, ketryx.KetryxAPIRepository{
			AuthToken:  repository.AuthToken.ValueString(),
			AuthUser:   repository.AuthUser.ValueString(),
			MainRef:    repository.MainRef.ValueString(),
			ReleaseRef: repository.ReleaseRef.ValueString(),
			Url:        repository.Url.ValueString(),
		})
	}

	// Update the project
	_, err = p.client.ProjectPost(ctx, request, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating project",
			"Could not update project, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state
	diags = resp.State.Set(ctx, plan)
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

type projectResourceModel struct {
	Id         types.String                     `tfsdk:"id"`
	Name       types.String                     `tfsdk:"name"`
	Repository []projectRepositoryResourceModel `tfsdk:"repository"`
	// TODO Figure out settings field
	// Settings   types.String                     `tfsdk:"settings"`
}

type projectRepositoryResourceModel struct {
	AuthToken  types.String `tfsdk:"auth_token"`
	AuthUser   types.String `tfsdk:"auth_user"`
	MainRef    types.String `tfsdk:"main_ref"`
	ReleaseRef types.String `tfsdk:"release_ref"`
	Url        types.String `tfsdk:"url"`
}
