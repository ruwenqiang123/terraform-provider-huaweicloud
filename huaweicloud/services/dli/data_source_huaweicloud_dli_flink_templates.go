// Generated by PMS #108
package dli

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDliFlinkTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDliFlinkTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the flink template to be queried.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the flink template to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the flink template to be queried.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All templates that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of template.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of template.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of template.`,
						},
						"sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The stream SQL statement.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of template.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the template.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the template.`,
						},
					},
				},
			},
		},
	}
}

type FlinkTemplatesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newFlinkTemplatesDSWrapper(d *schema.ResourceData, meta interface{}) *FlinkTemplatesDSWrapper {
	return &FlinkTemplatesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDliFlinkTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newFlinkTemplatesDSWrapper(d, meta)
	lisFliSqlJobTemRst, err := wrapper.ListFlinkSqlJobTemplates()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listFlinkSqlJobTemplatesToSchema(lisFliSqlJobTemRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DLI GET /v1.0/{project_id}/streaming/job-templates
func (w *FlinkTemplatesDSWrapper) ListFlinkSqlJobTemplates() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dli")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/streaming/job-templates"
	params := map[string]any{
		"name": w.Get("name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("template_list.templates", "offset", "limit", 50).
		Filter(
			filters.New().From("template_list.templates").
				Where("template_id", "=", w.GetToInt("template_id")).
				Where("job_type", "=", w.Get("type")),
		).
		Request().
		Result()
}

func (w *FlinkTemplatesDSWrapper) listFlinkSqlJobTemplatesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("templates", schemas.SliceToList(body.Get("template_list.templates"),
			func(template gjson.Result) any {
				return map[string]any{
					"id":          w.setTemLisTemTemId(template),
					"name":        template.Get("name").Value(),
					"type":        template.Get("job_type").Value(),
					"sql":         template.Get("sql_body").Value(),
					"description": template.Get("desc").Value(),
					"created_at":  w.setTemLisTemCreTim(template),
					"updated_at":  w.setTemLisTemUpdTim(template),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*FlinkTemplatesDSWrapper) setTemLisTemTemId(data gjson.Result) string {
	return data.Get("template_id").String()
}

func (*FlinkTemplatesDSWrapper) setTemLisTemCreTim(data gjson.Result) string {
	rawDate := data.Get("create_time").Int()
	return utils.FormatTimeStampRFC3339(rawDate/1000, false)
}

func (*FlinkTemplatesDSWrapper) setTemLisTemUpdTim(data gjson.Result) string {
	rawDate := data.Get("update_time").Int()
	return utils.FormatTimeStampRFC3339(rawDate/1000, false)
}
