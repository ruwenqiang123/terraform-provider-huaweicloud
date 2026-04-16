---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_catalog_metadata_task"
description: |-
  Manages DataArts Catalog metadata task resource within HuaweiCloud.
---

# huaweicloud_dataarts_catalog_metadata_task

Manages DataArts Catalog metadata task resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "connection_id" {}
variable "name" {}
variable "directory_id" {}
variable "schedule_config" {
  type = list(object({
    cron_expression = string
    max_time_out    = int
    interval        = string
    schedule_type   = string
    start_time      = string
    enabled         = int
  }))
}
variable "data_source_type" {}
variable "description" {}

data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = var.workspace_id
  connection_id = var.connection_id
}

resource "huaweicloud_dataarts_catalog_metadata_task" "test" {
  workspace_id                  = var.workspace_id
  name                          = var.name
  dir_id                        = var.directory_id

  schedule_config {
    cron_expression = "00 */15 9-23 * * ?"
    max_time_out    = 60
    end_time        = "2026-04-14T23:59:59 +08"
    interval        = "15 minutes"
    schedule_type   = "CRON"
    start_time      = "2026-04-14T00:00:00 +08"
    enabled         = 1
  } 

  data_source_type = var.data_source_type
  task_config      = jsonencode({
    data_connection_name        = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, "")
    data_connection_id          = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].id, "")
    data_connection_create_time = "1775892315978"
    databaseName                = [
      "tf_test_randx",
      "tpch"
    ]
    tableName                   = [
      "tf_test_randx.tf_test_randx",
      "tpch.part",
      "tpch.region",
      "tpch.customer"
    ]
    alive_object_policy         = "3"
    deleted_obkect_policy       = "1"
    deleted_object_policy       = "10"
    enableDataProfile           = true
    enableDataClassification    = false
    sampling                    = "10"
    queue                       = "default"
    unique                      = true
  })
  description      = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the metadata task is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the metadata task
  belongs.

* `name` - (Required, String) Specifies the name of the metadata task.

* `dir_id` - (Required, String) Specifies the directory name of the metadata task.

* `schedule_config` - (Required, List) Specifies the dispatch information of the metadata task.  
  The [schedule_config](#dataarts_catalog_metadata_task_schedule_config) structure is documented below.

* `data_source_type` - (Required, String, NonUpdatable) Specifies the data source type of the metadata task.

* `task_config` - (Required, String) Specifies the config information of the metadata task, in JSON format.

* `description` - (Optional, String) Specifies the description of the metadata task.

* `terminal_before_modify` - (Optional, Bool) Specifies whether to force terminal matadata task before update or delete
  it.

<a name="dataarts_catalog_metadata_task_schedule_config"></a>
The `schedule_config` block supports:

* `cron_expression` - (Optional, String) Specifies the cron expression of the schedule task.

* `end_time` - (Optional, String) Specifies the end time of the schedule task.

* `max_time_out` - (Optional, Int) Specifies the max time out of the schedule task.

* `interval` - (Optional, String) Specifies the interval time of the schedule task.

* `schedule_type` - (Required, String) Specifies the schedule type of the schedule task.

* `start_time` - (Optional, String) Specifies the start time of the schedule task.

* `enabled` - (Optional, Int) Specifies whether to enable the schedule task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `schedule_config` - The dispatch information of the metadata task.  
  The [schedule_config](#dataarts_catalog_metadata_task_schedule_config_attr) structure is documented below.

* `create_time` - The create time of the metadata task.

* `update_time` - The latest update time of the metadata task.

* `user_id` - The user ID which the metadata task is created.

* `user_name` - The user name which the metadata task is created.

* `path` - The directory path of the metadata task.

* `last_run_time` - The last run time of the metadata task.

* `start_time` - The start time of the metadata task.

* `end_time` - The end time of the metadata task.

* `next_run_time` - The next run time of the metadata task.

* `duty_person` - The duty person of the metadata task.

<a name="dataarts_catalog_metadata_task_schedule_config_attr"></a>
The `schedule_config` block supports:

* `job_id` - The job ID of the schedule task.

## Import

The catalog metadata task can be imported using `workspace_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_catalog_metadata_task.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `task_config`.
It is generally recommended running `terraform plan` after importing an metadata task.
You can then decide if changes should be applied to the metadata task, or the resource definition should be updated to
align with the metadata task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_catalog_metadata_task" "test" {
  ...

  lifecycle {
    ignore_changes = [
      task_config,
    ]
  }
}
```
