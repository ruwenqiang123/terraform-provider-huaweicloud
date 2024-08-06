// Generated by PMS #284
package organizations

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceOrganizationsTagPolicyServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsTagPolicyServicesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"services": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the services that support enforcement with tag policies.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the service name of the service.`,
						},
						"resource_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the resource types.`,
						},
						"support_all": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether resource_type support all services (wildcard *).`,
						},
					},
				},
			},
		},
	}
}

type TagPolicyServicesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTagPolicyServicesDSWrapper(d *schema.ResourceData, meta interface{}) *TagPolicyServicesDSWrapper {
	return &TagPolicyServicesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceOrganizationsTagPolicyServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTagPolicyServicesDSWrapper(d, meta)
	lisTagPolSerRst, err := wrapper.ListTagPolicyServices()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listTagPolicyServicesToSchema(lisTagPolSerRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API Organizations GET /v1/organizations/tag-policy-services
func (w *TagPolicyServicesDSWrapper) ListTagPolicyServices() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "organizations")
	if err != nil {
		return nil, err
	}

	uri := "/v1/organizations/tag-policy-services"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *TagPolicyServicesDSWrapper) listTagPolicyServicesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("services", schemas.SliceToList(body.Get("services"),
			func(services gjson.Result) any {
				return map[string]any{
					"service_name":   services.Get("service_name").Value(),
					"resource_types": schemas.SliceToStrList(services.Get("resource_types")),
					"support_all":    services.Get("support_all").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
