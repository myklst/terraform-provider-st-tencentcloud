package tencentcloud

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tencentCloudClbClient "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

var (
	_ datasource.DataSource              = &clbLoadBalancersDataSource{}
	_ datasource.DataSourceWithConfigure = &clbLoadBalancersDataSource{}
)

func NewClbDataSource() datasource.DataSource {
	return &clbLoadBalancersDataSource{}
}

type clbLoadBalancersDataSource struct {
	client *tencentCloudClbClient.Client
}

type clbLoadBalancersDataSourceModel struct {
	Id            types.String              `tfsdk:"id"`
	Name          types.String              `tfsdk:"name"`
	Tags          types.Map                 `tfsdk:"tags"`
	LoadBalancers []*clbLoadBalancersDetail `tfsdk:"load_balancers"`
}

type clbLoadBalancersDetail struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Tags types.Map    `tfsdk:"tags"`
}

func (d *clbLoadBalancersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_load_balancers"
}

func (d *clbLoadBalancersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides the Cloud Load Balancers of the current Tencent Cloud user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"tags": schema.MapAttribute{
				Description: "",
				ElementType: types.StringType,
				Optional:    true,
			},
			"load_balancers": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *clbLoadBalancersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(tencentCloudClients).clbClient
}

func (d *clbLoadBalancersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan *clbLoadBalancersDataSourceModel
	getPlanDiags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if getPlanDiags.HasError() {
		return
	}

	state := &clbLoadBalancersDataSourceModel{}

	state.Id = plan.Id
	state.Name = plan.Name
	state.Tags = plan.Tags
	state.LoadBalancers = []*clbLoadBalancersDetail{}

	// Create Describe Load Balancers Request
	describeLoadBalancersRequest := tencentCloudClbClient.NewDescribeLoadBalancersRequest()

	if !(plan.Name.IsUnknown() || plan.Name.IsNull()) {
		state.Name = plan.Name
		describeLoadBalancersRequest.LoadBalancerName = common.StringPtr(plan.Name.ValueString())
	}

	if !(plan.Id.IsUnknown() || plan.Id.IsNull()) {
		state.Id = plan.Id
		describeLoadBalancersRequest.LoadBalancerIds = []*string{common.StringPtr(state.Id.ValueString())}
	}

	describeLb := func() error {
		// Describe Load Balancers
		describeLoadBalancersResponse, err := d.client.DescribeLoadBalancers(describeLoadBalancersRequest)
		if err != nil {
			if terr, ok := err.(*errors.TencentCloudSDKError); ok {
				if isAbleToRetry(terr.GetCode()) {
					return err
				} else {
					return backoff.Permanent(err)
				}
			} else {
				return err
			}
		}

		// Read all Load Balancers
		for _, detail := range describeLoadBalancersResponse.Response.LoadBalancerSet {
			if len(detail.Tags) < 1 {
				clbDetail := &clbLoadBalancersDetail{
					Id:   types.StringValue(*detail.LoadBalancerId),
					Name: types.StringValue(*detail.LoadBalancerName),
					Tags: types.MapNull(types.StringType),
				}
				state.LoadBalancers = append(state.LoadBalancers, clbDetail)
				continue
			} else {

				// Initialize Tag Map
				clbTagMap := make(map[string]attr.Value)
				count := len(detail.Tags)
				for i := 0; i < count; i++ {
					clbTagMap[*detail.Tags[i].TagKey] = types.StringValue(*detail.Tags[i].TagValue)
				}

				// Determines whether Cloud Load Balancer is read
				clbOutput := true

				if !(plan.Tags.IsUnknown() || plan.Tags.IsNull()) {
					goInputMap := state.Tags.Elements()
					for inputKey, inputValue := range goInputMap {
						value, ok := clbTagMap[inputKey]
						if !ok || value != inputValue {
							clbOutput = false
							break
						}
					}
				}

				if clbOutput {
					clbDetail := &clbLoadBalancersDetail{
						Id:   types.StringValue(*detail.LoadBalancerId),
						Name: types.StringValue(*detail.LoadBalancerName),
						Tags: types.MapValueMust(types.StringType, clbTagMap),
					}
					state.LoadBalancers = append(state.LoadBalancers, clbDetail)
				}
			}
		}

		return nil
	}

	reconnectBackoff := backoff.NewExponentialBackOff()
	reconnectBackoff.MaxElapsedTime = 30 * time.Second

	err := backoff.Retry(describeLb, reconnectBackoff)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Describe Load Balancers",
			err.Error(),
		)
		return
	}

	setStateDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
