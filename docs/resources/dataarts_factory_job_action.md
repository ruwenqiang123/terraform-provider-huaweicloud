---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_job_action"
description: |-
  Use this resource to execute the DataArts Factory job action within HuaweiCloud.
---

# huaweicloud_dataarts_factory_job_action

Use this resource to execute the DataArts Factory job action within HuaweiCloud.

-> This resource is only a one-time action resource for changing job status. Deleting this resource will
   not change the current status, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "job_name" {}

resource "huaweicloud_dataarts_factory_job_action" "test" {
  workspace_id = var.workspace_id
  action       = "start"
  job_name     = var.job_name
  process_type = "BATCH"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the job is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `job_name` - (Required, String, NonUpdatable) Specified the name of the job to be performed.  
  The valid value is limited from `1` to `128` characters, only letters, numbers, hyphens (-), underscores (_),
  and periods (.) are allowed.

* `process_type` - (Required, String, NonUpdatable) Specified the type of the job to be performed.  
  The valid values are as follows:
  + **REAL_TIME**: Real-time processing.
  + **BATCH**: Batch processing.

* `action` - (Required, String) Specified the action type of the job to be performed.  
  The valid values are as follows:
  + **start**
  + **stop**

* `workspace_id` - (Optional, String, NonUpdatable) Specified the ID of the workspace to which the job belongs.
  If this parameter is not set, the default workspace is used by default.

* `job_params` - (Optional, List) Specifies the parameters of the job action.  
  The [job_params](#dataarts_factory_job_action_job_params) structure is documented below.  
  Only used when `action` is **start**.

* `start_date` - (Optional, Int) Specifies the start date of the job action when starting the job.  
  The format is **YYmmDD**, such as **20060102**.  
  Only used when `action` is **start**.

* `ignore_first_self_dep` - (Optional, Bool) Specifies whether to ignore the first self dependence when starting the
  job.  
  Only used when `action` is **start**.

<a name="dataarts_factory_job_action_job_params"></a>
The `job_params` block supports:

* `name` - (Required, String) Specifies the name of the job parameter.

* `value` - (Required, String) Specifies the value of the job parameter.

* `type` - (Optional, String) Specifies the type of the job parameter.  
  The valid values are as follows:
  + **variable**
  + **constants**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `job_name`.

* `status` - The current status of the job.
  + **NORMAL**
  + **STOPPED**
  + **SCHEDULING**
  + **PAUSED**
  + **EXCEPTION**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
