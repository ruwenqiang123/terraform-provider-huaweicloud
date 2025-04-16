// Generated by PMS #604
package tms

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

func DataSourceTmsTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsTagsRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the tag key.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the tag value.`,
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting field:`,
			},
			"order_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting method of the ` + "`" + `order_field` + "`" + ` field.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of tags.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the key of the tag.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the value of the tag.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when the tag is updated.`,
						},
					},
				},
			},
		},
	}
}

type TagsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTagsDSWrapper(d *schema.ResourceData, meta interface{}) *TagsDSWrapper {
	return &TagsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceTmsTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTagsDSWrapper(d, meta)
	listPredefineTagsRst, err := wrapper.ListPredefineTags()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listPredefineTagsToSchema(listPredefineTagsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API TMS GET /v1.0/predefine_tags
func (w *TagsDSWrapper) ListPredefineTags() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "tms")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/predefine_tags"
	params := map[string]any{
		"key":          w.Get("key"),
		"value":        w.Get("value"),
		"order_field":  w.Get("order_field"),
		"order_method": w.Get("order_method"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("tags", "marker", "marker").
		Request().
		Result()
}

func (w *TagsDSWrapper) listPredefineTagsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("tags", schemas.SliceToList(body.Get("tags"),
			func(tags gjson.Result) any {
				return map[string]any{
					"key":         tags.Get("key").Value(),
					"value":       tags.Get("value").Value(),
					"update_time": tags.Get("update_time").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
