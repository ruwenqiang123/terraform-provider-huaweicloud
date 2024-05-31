package rds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3.1/{project_id}/instances/recovery
// @API RDS GET /v3/{project_id}/instances
func ResourceRdsRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsRestoreCreate,
		ReadContext:   resourceRdsRestoreRead,
		DeleteContext: resourceRdsRestoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the target instance ID.`,
			},
			"source_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the source instance ID.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the restoration type.`,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the backup to be restored.`,
			},
			"restore_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"backup_id"},
				Description:  `Specifies the time point of data restoration in the UNIX timestamp format.`,
			},
			"database_name": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the databases that will be restored.`,
			},
		},
	}
}

func resourceRdsRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3.1/{project_id}/instances/recovery"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	targetInstanceId := d.Get("target_instance_id").(string)
	backupId := d.Get("backup_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreateRestoreBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Get("target_instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error restoring from backup(%s) to RDS instance (%s): %s", backupId, targetInstanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("unable to find the job_id from the response: %s", err)
	}

	d.SetId(jobId.(string))

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for restoring from backup(%s) to RDS instance(%s) to complete: %s",
			backupId, targetInstanceId, err)
	}

	return nil
}

func buildCreateRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	databaseName := d.Get("database_name").(map[string]interface{})
	bodyParams := map[string]interface{}{
		"source": map[string]interface{}{
			"instance_id":   d.Get("source_instance_id"),
			"type":          utils.ValueIngoreEmpty(d.Get("type")),
			"backup_id":     utils.ValueIngoreEmpty(d.Get("backup_id")),
			"restore_time":  utils.ValueIngoreEmpty(d.Get("restore_time")),
			"database_name": utils.ValueIngoreEmpty(databaseName),
		},
		"target": map[string]interface{}{
			"instance_id": d.Get("target_instance_id"),
		},
	}
	if len(databaseName) == 0 {
		bodyParams["source"].(map[string]interface{})["restore_all_database"] = true
	}
	return bodyParams
}

func resourceRdsRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restoration record is not supported. The restoration record is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before restoration."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
