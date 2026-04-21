---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_catalog_metadata_task_action"
description: |-
  Use this resource to operate a DataArts Catalog metadata task within HuaweiCloud.
---

# huaweicloud_dataarts_catalog_metadata_task_action

Use this resource to operate a DataArts Catalog metadata task within HuaweiCloud.

-> This resource is only a one-time action resource for operating a metadata task. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file,
   but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "task_id" {}

resource "huaweicloud_dataarts_catalog_metadata_task_action" "test" {
  workspace_id = var.workspace_id
  task_id      = var.task_id
  action       = "run"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the metadata task is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the metadata task
  belongs.

* `task_id` - (Required, String, NonUpdatable) Specifies the task ID of the metadata task.

* `action` - (Required, String, NonUpdatable) Specifies the supported action of the metadata task status.  
  The valid values are as follows:
  + **run**
  + **runimmediate**
  + **stop**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
