// Generated by PMS #78
package dns

import (
	"context"
	"strings"

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

func DataSourceFloatingPtrrecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFloatingPtrrecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"record_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the PTR record.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the EIP address of the PTR record.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the domain name of the PTR record.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID corresponding to the PTR record.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the key/value pairs to associate with the PTR record.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the PTR record.`,
			},
			"ptrrecords": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the PTR records.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the PTR record.`,
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The EIP address corresponding to the PTR record.`,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name of the PTR record.`,
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The valid cache time of the PTR record (in seconds).`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID corresponding to the PTR record.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs to associate with the PTR record.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the PTR record.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the PTR record.`,
						},
					},
				},
			},
		},
	}
}

type FloatingPtrrecordsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newFloatingPtrrecordsDSWrapper(d *schema.ResourceData, meta interface{}) *FloatingPtrrecordsDSWrapper {
	return &FloatingPtrrecordsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceFloatingPtrrecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newFloatingPtrrecordsDSWrapper(d, meta)
	lisPtrRecRst, err := wrapper.ListPtrRecords()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listPtrRecordsToSchema(lisPtrRecRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DNS GET /v2/reverse/floatingips
func (w *FloatingPtrrecordsDSWrapper) ListPtrRecords() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dns")
	if err != nil {
		return nil, err
	}

	uri := "/v2/reverse/floatingips"
	params := map[string]any{
		"enterprise_project_id": w.Get("enterprise_project_id"),
		"tags":                  w.getTag(),
		"status":                w.Get("status"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("floatingips", "offset", "limit", 0).
		Filter(
			filters.New().From("floatingips").
				Where("id", "=", w.Get("record_id")).
				Where("address", "=", w.Get("public_ip")).
				Where("ptrdname", "=", w.Get("domain_name")),
		).
		Request().
		Result()
}

func (w *FloatingPtrrecordsDSWrapper) listPtrRecordsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("ptrrecords", schemas.SliceToList(body.Get("floatingips"),
			func(ptrrecord gjson.Result) any {
				return map[string]any{
					"id":                    ptrrecord.Get("id").Value(),
					"public_ip":             ptrrecord.Get("address").Value(),
					"domain_name":           ptrrecord.Get("ptrdname").Value(),
					"ttl":                   ptrrecord.Get("ttl").Value(),
					"enterprise_project_id": ptrrecord.Get("enterprise_project_id").Value(),
					"tags":                  w.setFloatingipsTag(ptrrecord),
					"description":           ptrrecord.Get("description").Value(),
					"status":                ptrrecord.Get("status").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (w *FloatingPtrrecordsDSWrapper) getTag() string {
	raw := w.Get("tags")
	if raw == nil {
		return ""
	}

	tags := raw.(map[string]interface{})
	tagsList := make([]string, 0, len(tags))
	for k, v := range tags {
		tagsList = append(tagsList, k+","+v.(string))
	}
	return strings.Join(tagsList, "|")
}

func (*FloatingPtrrecordsDSWrapper) setFloatingipsTag(data gjson.Result) map[string]string {
	tags := make(map[string]string)
	tagList := data.Get("tags").Array()
	for _, v := range tagList {
		tags[v.Get("key").String()] = v.Get("value").String()
	}
	return tags
}
