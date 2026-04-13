---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_script_execute"
description: |-
  Use this resource to execute a DataArts Factory script once within HuaweiCloud.
---

# huaweicloud_dataarts_factory_script_execute

Use this resource to execute a DataArts Factory script once within HuaweiCloud.

-> This resource is only a one-time action resource for executing a script. Deleting this resource will not remove the
   execution record (only when the script is deleted can the record be cleared), but will only remove the resource
   information from the tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "script_name" {}
variable "execution_parameters" {
  type = map(string)
}

resource "huaweicloud_dataarts_factory_script_execute" "test" {
  workspace_id = var.workspace_id
  script_name  = var.script_name
  params       = jsonencode(var.execution_parameters)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workspace is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the script belongs.

* `script_name` - (Required, String, NonUpdatable) Specifies the name of the script to be executed.

* `params` - (Optional, String, NonUpdatable) Specifies the execution parameters passed to the script content, in JSON
  format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The execution status of the script instance.

* `message` - The message when the script instance execution fails.

## Timeouts

* `create` - Default is 10 minutes.
