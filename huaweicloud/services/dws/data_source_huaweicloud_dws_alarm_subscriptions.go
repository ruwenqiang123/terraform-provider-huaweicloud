// Generated by PMS #293
package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/alarm-subs
func DataSourceAlarmSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmSubscriptionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"subscriptions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All alarm subscriptions that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the alarm subscription.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the alarm subscription.`,
						},
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Whether alarm subscription is enabled.`,
						},
						"alarm_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The level of the alarm subscription.`,
						},
						"notification_target_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of notification topic corresponding to the alarm subscription.`,
						},
						"notification_target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of notification topic corresponding to the alarm subscription.`,
						},
						"notification_target_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of notification topic corresponding to the alarm subscription.`,
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time zone of the alarm subscription.`,
						},
						"language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The language of the alarm subscription.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlarmSubscriptionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	subscriptions, err := queryAlarmSubscriptions(client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("subscriptions", flattenAlarmSubscriptions(subscriptions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryAlarmSubscriptions(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/alarm-subs"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		// The offset indicates the page number.
		// The default value is 0, which represents the first page.
		listPathWithPage := fmt.Sprintf("%s?offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving alarm subscriptions: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		subscriptions := utils.PathSearch("alarm_subscriptions", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, subscriptions...)
		if len(result) == int(utils.PathSearch("count", respBody, float64(0)).(float64)) {
			break
		}
		offset++
	}
	return result, nil
}

func flattenAlarmSubscriptions(all []interface{}) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, v := range all {
		result[i] = map[string]interface{}{
			"id":                       utils.PathSearch("id", v, nil),
			"name":                     utils.PathSearch("name", v, nil),
			"enable":                   utils.PathSearch("enable", v, 0),
			"alarm_level":              utils.PathSearch("alarm_level", v, nil),
			"notification_target_type": utils.PathSearch("notification_target_type", v, nil),
			"notification_target":      utils.PathSearch("notification_target", v, nil),
			"notification_target_name": utils.PathSearch("notification_target_name", v, nil),
			"time_zone":                utils.PathSearch("time_zone", v, nil),
			"language":                 utils.PathSearch("language", v, nil),
		}
	}
	return result
}
