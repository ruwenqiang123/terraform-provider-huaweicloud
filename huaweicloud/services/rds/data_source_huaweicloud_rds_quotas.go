// Generated by PMS #560
package rds

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

func DataSourceRdsQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the objects in the quota list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the resource list objects.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"quota": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the project resource quota.`,
									},
									"used": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of used resources.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the project resource type. The value can be **instance**.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type QuotasDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newQuotasDSWrapper(d *schema.ResourceData, meta interface{}) *QuotasDSWrapper {
	return &QuotasDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRdsQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newQuotasDSWrapper(d, meta)
	showQuotasRst, err := wrapper.ShowQuotas()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showQuotasToSchema(showQuotasRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API RDS GET /v3/{project_id}/quotas
func (w *QuotasDSWrapper) ShowQuotas() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/quotas"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *QuotasDSWrapper) showQuotasToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("quotas", schemas.ObjectToList(body.Get("quotas"),
			func(quotas gjson.Result) any {
				return map[string]any{
					"resources": schemas.SliceToList(quotas.Get("resources"),
						func(resources gjson.Result) any {
							return map[string]any{
								"quota": resources.Get("quota").Value(),
								"used":  resources.Get("used").Value(),
								"type":  resources.Get("type").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
