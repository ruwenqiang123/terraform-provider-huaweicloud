---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_action"
description: |-
  Use this resource to operate DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_action

Use this resource to operate DWS cluster within HuaweiCloud.

-> This resource is only a one-time action resource for operating the DWS cluster. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "cluster_id" {}

resource "huaweicloud_dws_cluster_action" "restart" {
  cluster_id = var.cluster_id
  action     = "restart"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cluster is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to be operated.

* `action` - (Required, String, NonUpdatable) Specifies the action type of the operation.  
  The valid values are as follows:
  + **start**
  + **stop**
  + **restart**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
