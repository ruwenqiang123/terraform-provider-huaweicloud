package elb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var domainAddressnNonUpdatableParams = []string{"loadbalancer_id", "ip_address"}

// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips/batch-enable
// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips/batch-disable
// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips
func ResourceDomainAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainAddressCreate,
		ReadContext:   resourceDomainAddressRead,
		UpdateContext: resourceDomainAddressUpdate,
		DeleteContext: resourceDomainAddressDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainAddressImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(domainAddressnNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDomainAddressBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"ips": []map[string]interface{}{
			{
				"ip_address": d.Get("ip_address"),
			},
		},
	}

	return params
}

func resourceDomainAddressCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		loadbalancerId = d.Get("loadbalancer_id").(string)
		httpUrl        = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips/batch-enable"
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{loadbalancer_id}", loadbalancerId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDomainAddressBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error adding IP address to domain resolution: %s", err)
	}

	d.SetId(loadbalancerId)

	return resourceDomainAddressRead(ctx, d, meta)
}

func resourceDomainAddressRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		ipAddress = d.Get("ip_address").(string)
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	address, err := GetDomainAddressInfomation(client, d.Id(), ipAddress)
	if err != nil {
		// If the load balancer not exsit, the query API return code is `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving domain resolution IP address information")
	}

	ipEnable := utils.PathSearch("enable", address, nil).(bool)
	if !ipEnable {
		// If the IP address not add to the domain name, the query API return code is `200` and the attribute `enable` is false,
		// this situation also need to CheckDeletedDiag
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving domain resolution IP address information")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ip_address", utils.PathSearch("ip_address", address, nil)),
		d.Set("enable", utils.PathSearch("enable", address, nil)),
		d.Set("type", utils.PathSearch("type", address, nil)),
		d.Set("domain_name", utils.PathSearch("domain_name", address, nil)),
		d.Set("created_at", utils.PathSearch("created_at", address, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", address, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDomainAddressInfomation(client *golangsdk.ServiceClient, loadbalancerId, ipAddress string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips?limit={limit}"
		limit   = 1
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{loadbalancer_id}", loadbalancerId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		addresses := utils.PathSearch("ips", respBody, make([]interface{}, 0)).([]interface{})
		address := utils.PathSearch(fmt.Sprintf("[?ip_address=='%s']|[0]", ipAddress), addresses, nil)
		if address != nil {
			return address, nil
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceDomainAddressUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDomainAddressDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips/batch-disable"
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{loadbalancer_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDomainAddressBodyParams(d),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error removing IP address from domain resolution")
	}

	return nil
}

func resourceDomainAddressImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<id>/<ip_address>', but got '%s'",
			importedId)
	}

	d.SetId(parts[0])

	mErr := multierror.Append(nil,
		d.Set("ip_address", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
