package dataarts

import (
	"context"
	"fmt"
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

var (
	catalogMetadataTaskNonUpdatableParams = []string{
		"workspace_id",
		"data_source_type",
	}

	metadataTaskExecuteNotFoundCodes = []string{
		"DLG.0818", // The workspace or metadata task does not exist when query or delete metadata task.
		"DLG.2208", // The delete metadata task ID does not exist.
		"DLG.2253", // The query metadata task ID does not exist.
		"DLG.2206", // Deletion is not supported for soft-deleted instances that are in a scheduling state.
	}
)

// @API DataArtsStudio POST /v3/{project_id}/metadata/tasks/create
// @API DataArtsStudio GET /v3/{project_id}/metadata/tasks/{task_id}
// @API DataArtsStudio PUT /v3/{project_id}/metadata/tasks/{task_id}
// @API DataArtsStudio DELETE /v3/{project_id}/metadata/tasks/{task_id}
// @API DataArtsStudio POST /v3/{project_id}/metadata/tasks/{task_id}/action
func ResourceCatalogMetadataTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCatalogMetadataTaskCreate,
		ReadContext:   resourceCatalogMetadataTaskRead,
		UpdateContext: resourceCatalogMetadataTaskUpdate,
		DeleteContext: resourceCatalogMetadataTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(catalogMetadataTaskNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceCatalogMetadataTaskImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the metadata task is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the metadata task belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the metadata task.",
			},
			"dir_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The directory name of the metadata task.",
			},
			"schedule_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        metadataTaskScheduleConfigElemSchema(),
				Description: "The dispatch information of the metadata task.",
			},
			"data_source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The data source type of the metadata task.",
			},
			"task_config": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The config information of the metadata task, in JSON format.",
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the metadata task.",
			},
			"terminal_before_modify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to force terminal matadata task before update or delete it.",
			},

			// Attributes.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the metadata task.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the metadata task.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID which the metadata task is created.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name which the metadata task is created.",
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The directory path of the metadata task.",
			},
			"last_run_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last run time of the metadata task.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the metadata task.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The end time of the metadata task.",
			},
			"next_run_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The next run time of the metadata task.",
			},
			"duty_person": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The duty person of the metadata task.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func metadataTaskScheduleConfigElemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cron_expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cron expression of the schedule task.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The end time of the schedule task.",
			},
			"max_time_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The max time out of the schedule task.",
			},
			"interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The interval time of the schedule task.",
			},
			"schedule_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The schedule type of the schedule task.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The start time of the schedule task.",
			},
			"enabled": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable the schedule task.",
			},
			"job_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The job ID of the schedule task.",
			},
		},
	}
	return &sc
}

func bulidCatalogMoreHeaders(workspaceId string) map[string]string {
	result := map[string]string{
		"Content-Type": "application/json",
	}
	if workspaceId != "" {
		result["workspace"] = workspaceId
	}
	return result
}

func buildMetadataTaskScheduleConfigBodyParams(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return nil
	}
	// The default integer triggers the structure change.
	return map[string]interface{}{
		"cron_expression": utils.ValueIgnoreEmpty(utils.PathSearch("cron_expression", params[0], "").(string)),
		"end_time":        utils.ValueIgnoreEmpty(utils.PathSearch("end_time", params[0], "").(string)),
		"max_time_out":    utils.PathSearch("max_time_out", params[0], 0).(int),
		"interval":        utils.ValueIgnoreEmpty(utils.PathSearch("interval", params[0], "").(string)),
		"schedule_type":   utils.ValueIgnoreEmpty(utils.PathSearch("schedule_type", params[0], "").(string)),
		"start_time":      utils.ValueIgnoreEmpty(utils.PathSearch("start_time", params[0], "").(string)),
		"enabled":         utils.PathSearch("enabled", params[0], 0).(int),
	}
}

func buildCreateCatalogMetadataTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name"),
		"description":      d.Get("description"),
		"dir_id":           d.Get("dir_id"),
		"schedule_config":  utils.RemoveNil(buildMetadataTaskScheduleConfigBodyParams(d.Get("schedule_config").([]interface{}))),
		"data_source_type": d.Get("data_source_type"),
		"task_config":      utils.StringToJson(d.Get("task_config").(string)),
	}
}

func createCatalogMetadataTask(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl     = "v3/{project_id}/metadata/tasks/create"
		workspaceId = d.Get("workspace_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createCatalogMetadataTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      bulidCatalogMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildCreateCatalogMetadataTaskBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createCatalogMetadataTaskOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceCatalogMetadataTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := createCatalogMetadataTask(client, d)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Catalog metadata task: %s", err)
	}

	taskId := utils.PathSearch("task_id", respBody, "").(string)
	if taskId == "" {
		return diag.Errorf("unable to find the ID of the DataArts Studio Catalog metadata task from the API response")
	}
	d.SetId(taskId)

	return resourceCatalogMetadataTaskRead(ctx, d, meta)
}

func GetCatalogMetadataTaskById(client *golangsdk.ServiceClient, workspaceId, taskId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/metadata/tasks/{task_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)
	getCatalogMetadataTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      bulidCatalogMoreHeaders(workspaceId),
	}

	resp, err := client.Request("GET", getPath, &getCatalogMetadataTaskOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func flattenScheduleConfig(scheduleConfig map[string]interface{}) []map[string]interface{} {
	if len(scheduleConfig) == 0 {
		return nil
	}

	enabled := int(utils.PathSearch("enabled", scheduleConfig, float64(0)).(float64))
	// The status of "3 in scheduling" and "5 paused" is uniformly changed to "scheduling status 1."
	if enabled == 3 || enabled == 5 {
		enabled = 1
	} else {
		// All other states are in the unscheduled state 0
		enabled = 0
	}
	return []map[string]interface{}{
		{
			"cron_expression": utils.PathSearch("cron_expression", scheduleConfig, nil),
			"end_time":        utils.PathSearch("end_time", scheduleConfig, nil),
			"max_time_out":    int(utils.PathSearch("max_time_out", scheduleConfig, float64(0)).(float64)),
			"interval":        utils.PathSearch("interval", scheduleConfig, nil),
			"schedule_type":   utils.PathSearch("schedule_type", scheduleConfig, nil),
			"start_time":      utils.PathSearch("start_time", scheduleConfig, nil),
			"job_id":          int(utils.PathSearch("job_id", scheduleConfig, float64(0)).(float64)),
			"enabled":         enabled,
		},
	}
}

func resourceCatalogMetadataTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		taskId      = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	resp, err := GetCatalogMetadataTaskById(client, workspaceId, taskId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors[0].error_code", metadataTaskExecuteNotFoundCodes...),
			fmt.Sprintf("error retrieving DataArts Catalog metadata task (%s)", taskId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", resp, nil)),
		d.Set("description", utils.PathSearch("description", resp, nil)),
		d.Set("user_id", utils.PathSearch("user_id", resp, nil)),
		d.Set("dir_id", utils.PathSearch("dir_id", resp, nil)),
		d.Set("schedule_config", flattenScheduleConfig(utils.PathSearch("schedule_config", resp,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("data_source_type", utils.PathSearch("data_source_type", resp, nil)),
		d.Set("task_config", utils.JsonToString(utils.PathSearch("task_config", resp, nil))),
		d.Set("create_time", utils.PathSearch("create_time", resp, nil)),
		d.Set("update_time", utils.PathSearch("update_time", resp, nil)),
		d.Set("user_name", utils.PathSearch("user_name", resp, nil)),
		d.Set("path", utils.PathSearch("path", resp, nil)),
		d.Set("last_run_time", utils.PathSearch("last_run_time", resp, nil)),
		d.Set("start_time", utils.PathSearch("start_time", resp, nil)),
		d.Set("end_time", utils.PathSearch("end_time", resp, nil)),
		d.Set("next_run_time", utils.PathSearch("next_run_time", resp, nil)),
		d.Set("duty_person", utils.PathSearch("duty_person", resp, nil)),
		d.Set("terminal_before_modify", d.Get("terminal_before_modify")),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateCatalogMetadataTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name"),
		"description":      d.Get("description"),
		"dir_id":           d.Get("dir_id"),
		"schedule_config":  utils.RemoveNil(buildMetadataTaskScheduleConfigBodyParams(d.Get("schedule_config").([]interface{}))),
		"data_source_type": d.Get("data_source_type"),
		"task_config":      utils.StringToJson(d.Get("task_config").(string)),
	}
}

func stopCatalogMetadataTask(client *golangsdk.ServiceClient, workspaceId, taskId string) error {
	httpUrl := "v3/{project_id}/metadata/tasks/{task_id}/action?action=stop"

	stopPath := client.Endpoint + httpUrl
	stopPath = strings.ReplaceAll(stopPath, "{project_id}", client.ProjectID)
	stopPath = strings.ReplaceAll(stopPath, "{task_id}", taskId)

	stopCatalogMetadataTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      bulidCatalogMoreHeaders(workspaceId),
	}

	_, err := client.Request("POST", stopPath, &stopCatalogMetadataTaskOpt)
	return err
}

func updateCatalogMetadataTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/metadata/tasks/{task_id}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{task_id}", d.Id())

	updateCatalogMetadataTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      bulidCatalogMoreHeaders(d.Get("workspace_id").(string)),
		JSONBody:         utils.RemoveNil(buildUpdateCatalogMetadataTaskBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateCatalogMetadataTaskOpt)
	return err
}

func resourceCatalogMetadataTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		taskId      = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	if d.HasChangeExcept("terminal_before_modify") {
		if terminalBeforeModify, ok := d.GetOk("terminal_before_modify"); ok && terminalBeforeModify.(bool) {
			scheduleConfigs := d.Get("schedule_config").([]interface{})
			if len(scheduleConfigs) > 0 && utils.PathSearch("enabled", scheduleConfigs[0], 0).(int) == 1 {
				err = stopCatalogMetadataTask(client, workspaceId, taskId)
				if err != nil {
					return diag.Errorf("error stoping metadata task (%s): %s", taskId, err)
				}
			}
		}

		err = updateCatalogMetadataTask(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCatalogMetadataTaskRead(ctx, d, meta)
}

func deleteCatalogMetadataTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/metadata/tasks/{task_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", d.Id())

	deleteCatalogMetadataTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      bulidCatalogMoreHeaders(d.Get("workspace_id").(string)),
		OkCodes:          []int{200},
	}

	_, err := client.Request("DELETE", deletePath, &deleteCatalogMetadataTaskOpt)
	return err
}

func resourceCatalogMetadataTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		taskId      = d.Id()
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	if terminalBeforeModify, ok := d.GetOk("terminal_before_modify"); ok && terminalBeforeModify.(bool) {
		scheduleConfigs := d.Get("schedule_config").([]interface{})
		if len(scheduleConfigs) > 0 && utils.PathSearch("enabled", scheduleConfigs[0], 0).(int) == 1 {
			err = stopCatalogMetadataTask(client, workspaceId, taskId)
			if err != nil {
				return diag.Errorf("error stoping metadata task (%s): %s", taskId, err)
			}
		}
	}

	err = deleteCatalogMetadataTask(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors[0].error_code", metadataTaskExecuteNotFoundCodes...),
			fmt.Sprintf("error deleting DataArts Catalog metadata task (%s)", taskId),
		)
	}

	return nil
}

func resourceCatalogMetadataTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', "+
			"but got '%s'", importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
