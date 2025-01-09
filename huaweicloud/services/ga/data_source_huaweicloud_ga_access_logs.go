// Generated by PMS #503
package ga

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceGaAccessLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaAccessLogsRead,

		Schema: map[string]*schema.Schema{
			"log_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the access log.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the access log.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the resource to which the access log belongs.`,
			},
			"resource_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the ID list of the resource to which the access log belongs.`,
			},
			"logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the access logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the access log.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the access log.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource to which the access log belongs.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource to which the access log belongs.`,
						},
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group to which the access log belongs.`,
						},
						"log_stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log stream to which the access log belongs.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the access log, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the access log, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

type AccessLogsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAccessLogsDSWrapper(d *schema.ResourceData, meta interface{}) *AccessLogsDSWrapper {
	return &AccessLogsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceGaAccessLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAccessLogsDSWrapper(d, meta)
	listLogtanksRst, err := wrapper.ListLogtanks()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listLogtanksToSchema(listLogtanksRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API GA GET /v1/logtanks
func (w *AccessLogsDSWrapper) ListLogtanks() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ga")
	if err != nil {
		return nil, err
	}

	uri := "/v1/logtanks"
	params := map[string]any{
		"id":            w.Get("log_id"),
		"status":        w.Get("status"),
		"resource_ids":  w.ListToArray("resource_ids"),
		"resource_type": w.Get("resource_type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("logtanks", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *AccessLogsDSWrapper) listLogtanksToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("logs", schemas.SliceToList(body.Get("logtanks"),
			func(logs gjson.Result) any {
				return map[string]any{
					"id":            logs.Get("id").Value(),
					"status":        logs.Get("status").Value(),
					"resource_type": logs.Get("resource_type").Value(),
					"resource_id":   logs.Get("resource_id").Value(),
					"log_group_id":  logs.Get("log_group_id").Value(),
					"log_stream_id": logs.Get("log_stream_id").Value(),
					"created_at":    w.setLogtanksCreatedAt(logs),
					"updated_at":    w.setLogtanksUpdatedAt(logs),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*AccessLogsDSWrapper) setLogtanksCreatedAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("created_at").String())/1000, false)
}

func (*AccessLogsDSWrapper) setLogtanksUpdatedAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("updated_at").String())/1000, false)
}
