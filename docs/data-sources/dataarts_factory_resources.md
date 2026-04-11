---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_resources"
description: |-
  Use this data source to get the list of DataArts Factory resources within HuaweiCloud.
---
# huaweicloud_dataarts_factory_resources

Use this data source to get the list of DataArts Factory resources within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_factory_resources" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the resources are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Optional, String) Specifies the ID of the workspace to which the resources belong.  
  If omitted, the default workspace will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of resources that match the filter parameters.  
  The [resources](#factory_resources) structure is documented below.

<a name="factory_resources"></a>
The `resources` block supports:

* `name` - The name of the resource.

* `type` - The type of the resource.

* `location` - The OBS path of the resource file.

* `directory` - The directory where the resource is located.

* `description` - The description of the resource.

* `depend_files` - The list of dependent file paths.
