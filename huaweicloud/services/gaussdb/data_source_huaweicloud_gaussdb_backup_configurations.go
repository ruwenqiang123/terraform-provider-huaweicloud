package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/backups/config
func DataSourceBackupConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rate_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"file_split_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"prefetch_block": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_standby_backup": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"close_compression": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_backup_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_backup_media_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_parallel_degree": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceBackupConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/backups/config"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the GaussDB backup configurations: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("rate_limit", utils.PathSearch("rate_limit", getRespBody, nil)),
		d.Set("file_split_size", utils.PathSearch("file_split_size", getRespBody, nil)),
		d.Set("prefetch_block", utils.PathSearch("prefetch_block", getRespBody, nil)),
		d.Set("enable_standby_backup", utils.PathSearch("enable_standby_backup", getRespBody, nil)),
		d.Set("close_compression", utils.PathSearch("close_compression", getRespBody, nil)),
		d.Set("default_backup_method", utils.PathSearch("default_backup_method", getRespBody, nil)),
		d.Set("default_backup_media_type", utils.PathSearch("default_backup_media_type", getRespBody, nil)),
		d.Set("backup_parallel_degree", utils.PathSearch("backup_parallel_degree", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
