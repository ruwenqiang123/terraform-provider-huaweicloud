// Generated by PMS #208
package dns

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

func DataSourceDNSCustomLines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSCustomLinesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"line_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the custom line. Fuzzy search is supported.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the custom line. Fuzzy search is supported.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The IP address used to query custom line which is in the IP address range.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the custom line.`,
			},
			"lines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All custom lines that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the custom line.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the custom line.`,
						},
						"ip_segments": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The IP address range of the custom line.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the custom line.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the custom line.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the custom line, in RFC339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the custom line, in RFC339 format.`,
						},
					},
				},
			},
		},
	}
}

type CustomLinesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCustomLinesDSWrapper(d *schema.ResourceData, meta interface{}) *CustomLinesDSWrapper {
	return &CustomLinesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDNSCustomLinesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCustomLinesDSWrapper(d, meta)
	listCustomLineRst, err := wrapper.ListCustomLine()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listCustomLineToSchema(listCustomLineRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DNS GET /v2.1/customlines
func (w *CustomLinesDSWrapper) ListCustomLine() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dns_region")
	if err != nil {
		return nil, err
	}

	uri := "/v2.1/customlines"
	params := map[string]any{
		"line_id": w.Get("line_id"),
		"name":    w.Get("name"),
		"status":  w.Get("status"),
		"ip":      w.Get("ip"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("lines", "offset", "limit", 0).
		Request().
		Result()
}

func (w *CustomLinesDSWrapper) listCustomLineToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("lines", schemas.SliceToList(body.Get("lines"),
			func(lines gjson.Result) any {
				return map[string]any{
					"id":          lines.Get("line_id").Value(),
					"name":        lines.Get("name").Value(),
					"ip_segments": schemas.SliceToStrList(lines.Get("ip_segments")),
					"status":      lines.Get("status").Value(),
					"description": lines.Get("description").Value(),
					"created_at":  w.setLinesCreatedAt(lines),
					"updated_at":  w.setLinesUpdatedAt(lines),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*CustomLinesDSWrapper) setLinesCreatedAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("created_at").String(), "2006-01-02T15:04:05")/1000, false)
}

func (*CustomLinesDSWrapper) setLinesUpdatedAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("updated_at").String(), "2006-01-02T15:04:05")/1000, false)
}
