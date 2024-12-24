package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live PUT /v1/{project_id}/notifications/publish
// @API Live GET /v1/{project_id}/notifications/publish
// @API Live DELETE /v1/{project_id}/notifications/publish
func ResourceNotificationCofiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCofigCreate,
		ReadContext:   resourceNotificationCofigRead,
		UpdateContext: resourceNotificationCofigUpdate,
		DeleteContext: resourceNotificationCofigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNotificationCofigImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ingest domain name to which the notification cofiguration belongs.`,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the callback URL.`,
			},
			// This parameter can be left blank.
			"auth_sign_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Specifies the authentication key.`,
			},
			// This parameter can be left blank.
			"call_back_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the area where the server for receiving callback notification is located.`,
			},
		},
	}
}

func resourceNotificationCofigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = createOrUpdateNotificationCofig(client, d)
	if err != nil {
		return diag.Errorf("error creating notification cofiguration: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceNotificationCofigRead(ctx, d, meta)
}

func createOrUpdateNotificationCofig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	notificationHttpUrl := "v1/{project_id}/notifications/publish"
	notificationPath := client.Endpoint + notificationHttpUrl
	notificationPath = strings.ReplaceAll(notificationPath, "{project_id}", client.ProjectID)
	notificationPath = fmt.Sprintf("%s?domain=%v", notificationPath, d.Get("domain_name"))

	notificationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildNotificationCofigBodyParams(d),
	}

	_, err := client.Request("PUT", notificationPath, &notificationOpt)
	return err
}

func buildNotificationCofigBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"url":            d.Get("url"),
		"auth_sign_key":  d.Get("auth_sign_key"),
		"call_back_area": d.Get("call_back_area"),
	}

	return params
}

func resourceNotificationCofigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
		getHttpUrl = "v1/{project_id}/notifications/publish"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%s", getPath, domainName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode),
			"error retrieving notification cofiguration")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	callbackUrl := utils.PathSearch("url", getRespBody, "").(string)
	if callbackUrl == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "notification cofiguration")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", domainName),
		d.Set("url", utils.PathSearch("url", getRespBody, nil)),
		d.Set("auth_sign_key", utils.PathSearch("auth_sign_key", getRespBody, nil)),
		d.Set("call_back_area", utils.PathSearch("call_back_area", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNotificationCofigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = createOrUpdateNotificationCofig(client, d)
	if err != nil {
		return diag.Errorf("error updating notification cofiguration: %s", err)
	}

	return resourceNotificationCofigRead(ctx, d, meta)
}

func resourceNotificationCofigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		deleteHttpUrl = "v1/{project_id}/notifications/publish"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?domain=%v", deletePath, d.Get("domain_name"))
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode),
			"error deleting notification cofiguration")
	}

	return nil
}

func resourceNotificationCofigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	if importedId == "" {
		return nil, fmt.Errorf("invalid format specified for import ID, 'domain_name' is empty")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("domain_name", importedId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
