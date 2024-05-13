// Generated by PMS #116
package dcs

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
)

func DataSourceDcsAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsAccountsRead,

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
				Description: `Specifies the instance ID.`,
			},
			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the account name.`,
			},
			"account_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the account type. The value can be **normal** or **default**.`,
			},
			"account_role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the account role. The value can be **read** or **write**.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the account description.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the account status.`,
			},
			"accounts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `ACL account list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Account ID.`,
						},
						"account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Account name.`,
						},
						"account_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Account type.`,
						},
						"account_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Account permissions.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Account description.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Account status.`,
						},
					},
				},
			},
		},
	}
}

type AccountsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAccountsDSWrapper(d *schema.ResourceData, meta interface{}) *AccountsDSWrapper {
	return &AccountsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDcsAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAccountsDSWrapper(d, meta)
	lisAclAccRst, err := wrapper.ListAclAccounts()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAclAccountsToSchema(lisAclAccRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DCS GET /v2/{project_id}/instances/{instance_id}/accounts
func (w *AccountsDSWrapper) ListAclAccounts() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dcs")
	if err != nil {
		return nil, err
	}

	d := w.ResourceData
	uri := "/v2/{project_id}/instances/{instance_id}/accounts"
	uri = strings.ReplaceAll(uri, "{instance_id}", d.Get("instance_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Filter(
			filters.New().From("accounts").
				Where("account_name", "=", w.Get("account_name")).
				Where("account_type", "=", w.Get("account_type")).
				Where("account_role", "=", w.Get("account_role")).
				Where("description", "=", w.Get("description")).
				Where("status", "=", w.Get("status")),
		).
		Request().
		Result()
}

func (w *AccountsDSWrapper) listAclAccountsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("accounts", schemas.SliceToList(body.Get("accounts"),
			func(account gjson.Result) any {
				return map[string]any{
					"id":           account.Get("account_id").Value(),
					"account_name": account.Get("account_name").Value(),
					"account_type": account.Get("account_type").Value(),
					"account_role": account.Get("account_role").Value(),
					"description":  account.Get("description").Value(),
					"status":       account.Get("status").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
