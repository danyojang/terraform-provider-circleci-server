package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &CircleCIServerProvider{}

type CircleCIServerProvider struct {
	version string
}

type CircleCIServerProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func (p *CircleCIServerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "circleci-server"
	resp.Version = p.version
}

func (p *CircleCIServerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "CircleCI Server host (e.g., https://cci.anduril.dev)",
				Required:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "CircleCI API token",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *CircleCIServerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CircleCIServerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := &CircleCIClient{
		Host:  data.Host.ValueString(),
		Token: data.Token.ValueString(),
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CircleCIServerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectFollowResource,
	}
}

func (p *CircleCIServerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CircleCIServerProvider{
			version: version,
		}
	}
}

type CircleCIClient struct {
	Host  string
	Token string
}
