// Generated by PMS #123
package waf

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

func DataSourceWafRulesBlacklist() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafRulesBlacklistRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the policy to which the blacklist and whitelist rules belong.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the blacklist or whitelist rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the blacklist or whitelist rule.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the blacklist or whitelist rule.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID to which the protection policies belong.`,
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the protective action of the blacklist and whitelist rule.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the blacklist and whitelist rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the blacklist or whitelist rule.`,
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the policy to which the blacklist and whitelist rule belongs.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the blacklist or whitelist rule.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the blacklist or whitelist rule.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the blacklist or whitelist rule.`,
						},
						"action": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The protective action of the blacklist and whitelist rule.`,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address included in the blacklist and whitelist rule.`,
						},
						"address_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The IP address group included in the blacklist and whitelist rule.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the IP address group.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the IP address group.`,
									},
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The number of IP addresses or IP address ranges in the IP address group.`,
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the blacklist and whitelist rule.`,
						},
					},
				},
			},
		},
	}
}

type RulesBlacklistDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newRulesBlacklistDSWrapper(d *schema.ResourceData, meta interface{}) *RulesBlacklistDSWrapper {
	return &RulesBlacklistDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceWafRulesBlacklistRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newRulesBlacklistDSWrapper(d, meta)
	lisWhiRulRst, err := wrapper.ListWhiteblackipRule()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listWhiteblackipRuleToSchema(lisWhiRulRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/whiteblackip
func (w *RulesBlacklistDSWrapper) ListWhiteblackipRule() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "waf")
	if err != nil {
		return nil, err
	}

	d := w.ResourceData
	uri := "/v1/{project_id}/waf/policy/{policy_id}/whiteblackip"
	uri = strings.ReplaceAll(uri, "{policy_id}", d.Get("policy_id").(string))
	params := map[string]any{
		"enterprise_project_id": w.Get("enterprise_project_id"),
		"name":                  w.Get("name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		PageSizePager("items", "page", "pagesize", 0).
		Filter(
			filters.New().From("items").
				Where("id", "=", w.Get("rule_id")).
				Where("status", "=", w.GetToInt("status")).
				Where("white", "=", w.GetToInt("action")),
		).
		Request().
		Result()
}

func (w *RulesBlacklistDSWrapper) listWhiteblackipRuleToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("rules", schemas.SliceToList(body.Get("items"),
			func(rule gjson.Result) any {
				return map[string]any{
					"id":          rule.Get("id").Value(),
					"policy_id":   rule.Get("policyid").Value(),
					"name":        rule.Get("name").Value(),
					"status":      rule.Get("status").Value(),
					"description": rule.Get("description").Value(),
					"action":      rule.Get("white").Value(),
					"ip_address":  rule.Get("addr").Value(),
					"address_group": schemas.SliceToList(rule.Get("ip_group"),
						func(addressGroup gjson.Result) any {
							return map[string]any{
								"id":   addressGroup.Get("id").Value(),
								"name": addressGroup.Get("name").Value(),
								"size": addressGroup.Get("size").Value(),
							}
						},
					),
					"created_at": w.setItemsTimestamp(rule),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*RulesBlacklistDSWrapper) setItemsTimestamp(data gjson.Result) string {
	rawDate := data.Get("timestamp").Int()
	return utils.FormatTimeStampRFC3339(rawDate/1000, false)
}
