// Generated by PMS #353
package dbss

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceDbssAuditRuleScopes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDbssAuditRuleScopesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the audit instance ID to which the audit scopes belong.`,
			},
			"scopes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the audit scopes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the audit scope.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the audit scope.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the audit scope.`,
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The action of the audit scope.`,
						},
						"exception_ips": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The exception IP addresses of the audit scope.`,
						},
						"source_ips": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source IP addresses of the audit scope.`,
						},
						"source_ports": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source ports of the audit scope.`,
						},
						"db_ids": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The database IDs associated with the audit scope.`,
						},
						"db_names": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The database names associated with the audit scope.`,
						},
						"db_users": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The database accounts associated with the audit scope.`,
						},
						"all_audit": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether is full audit.`,
						},
					},
				},
			},
		},
	}
}

type AuditRuleScopesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAuditRuleScopesDSWrapper(d *schema.ResourceData, meta interface{}) *AuditRuleScopesDSWrapper {
	return &AuditRuleScopesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDbssAuditRuleScopesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAuditRuleScopesDSWrapper(d, meta)
	lisAudRulScoRst, err := wrapper.ListAuditRuleScopes()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAuditRuleScopesToSchema(lisAudRulScoRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DBSS GET /v1/{project_id}/{instance_id}/dbss/audit/rule/scopes
func (w *AuditRuleScopesDSWrapper) ListAuditRuleScopes() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dbss")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/{instance_id}/dbss/audit/rule/scopes"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("scopes", "offset", "limit", 100).
		Request().
		Result()
}

func (w *AuditRuleScopesDSWrapper) listAuditRuleScopesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("scopes", schemas.SliceToList(body.Get("scopes"),
			func(scopes gjson.Result) any {
				return map[string]any{
					"id":            scopes.Get("id").Value(),
					"name":          scopes.Get("name").Value(),
					"status":        scopes.Get("status").Value(),
					"action":        scopes.Get("action").Value(),
					"exception_ips": scopes.Get("exception_ips").Value(),
					"source_ips":    scopes.Get("source_ips").Value(),
					"source_ports":  scopes.Get("source_ports").Value(),
					"db_ids":        scopes.Get("db_ids").Value(),
					"db_names":      scopes.Get("db_names").Value(),
					"db_users":      scopes.Get("db_users").Value(),
					"all_audit":     scopes.Get("all_audit").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
