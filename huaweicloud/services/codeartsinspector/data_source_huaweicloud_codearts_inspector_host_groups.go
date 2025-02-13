// Generated by PMS #542
package codeartsinspector

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

func DataSourceCodeartsInspectorHostGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsInspectorHostGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Specifies the group list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the group ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the group name.`,
						},
					},
				},
			},
		},
	}
}

type HostGroupsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newHostGroupsDSWrapper(d *schema.ResourceData, meta interface{}) *HostGroupsDSWrapper {
	return &HostGroupsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCodeartsInspectorHostGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newHostGroupsDSWrapper(d, meta)
	listGroupsRst, err := wrapper.ListGroups()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listGroupsToSchema(listGroupsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API VSS GET /v3/{project_id}/hostscan/groups
func (w *HostGroupsDSWrapper) ListGroups() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "vss")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/hostscan/groups"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("items", "offset", "limit", 0).
		Request().
		Result()
}

func (w *HostGroupsDSWrapper) listGroupsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("groups", schemas.SliceToList(body.Get("items"),
			func(groups gjson.Result) any {
				return map[string]any{
					"id":   groups.Get("id").Value(),
					"name": groups.Get("name").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
