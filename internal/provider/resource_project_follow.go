package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ProjectFollowResource{}
var _ resource.ResourceWithImportState = &ProjectFollowResource{}

func NewProjectFollowResource() resource.Resource {
	return &ProjectFollowResource{}
}

type ProjectFollowResource struct {
	client *CircleCIClient
}

type ProjectFollowResourceModel struct {
	ID           types.String `tfsdk:"id"`
	VCSType      types.String `tfsdk:"vcs_type"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
}

func (r *ProjectFollowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_follow"
}

func (r *ProjectFollowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Follows a CircleCI project using the v1.1 API",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vcs_type": schema.StringAttribute{
				MarkdownDescription: "VCS type (e.g., 'github', 'bitbucket')",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "Organization/username",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "Project/repository name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ProjectFollowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CircleCIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *CircleCIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ProjectFollowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectFollowResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call v1.1 follow API
	url := fmt.Sprintf("%s/api/v1.1/project/%s/%s/%s/follow",
		r.client.Host,
		data.VCSType.ValueString(),
		data.Organization.ValueString(),
		data.Project.ValueString())

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create request: %s", err))
		return
	}

	httpReq.Header.Set("Circle-Token", r.client.Token)
	httpReq.Header.Set("Accept", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to follow project: %s", err))
		return
	}
	defer httpResp.Body.Close()

	body, _ := io.ReadAll(httpResp.Body)

	if httpResp.StatusCode != 200 && httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Failed to follow project (status %d): %s", httpResp.StatusCode, string(body)),
		)
		return
	}

	// Set ID as composite key
	data.ID = types.StringValue(fmt.Sprintf("%s/%s/%s",
		data.VCSType.ValueString(),
		data.Organization.ValueString(),
		data.Project.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectFollowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectFollowResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// v1.1 API doesn't have a good way to check if following, so we just assume it exists
	// Real implementation could call v2 API to check project exists
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectFollowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectFollowResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Since all attributes require replace, update should never be called
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectFollowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectFollowResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call v1.1 unfollow API
	url := fmt.Sprintf("%s/api/v1.1/project/%s/%s/%s/unfollow",
		r.client.Host,
		data.VCSType.ValueString(),
		data.Organization.ValueString(),
		data.Project.ValueString())

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create request: %s", err))
		return
	}

	httpReq.Header.Set("Circle-Token", r.client.Token)
	httpReq.Header.Set("Accept", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to unfollow project: %s", err))
		return
	}
	defer httpResp.Body.Close()

	// Ignore errors on delete - project might already be unfollowed
}

func (r *ProjectFollowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: vcs_type/organization/project
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected format: vcs_type/organization/project, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vcs_type"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project"), parts[2])...)
}
