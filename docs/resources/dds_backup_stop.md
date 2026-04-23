---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_backup_stop"
description: |-
  Manages a resource to stop DDS backup within HuaweiCloud.
---

# huaweicloud_dds_backup_stop

Manages a resource to stop DDS backup within HuaweiCloud.

## Example Usage

```hcl
variable "backup_id" {}
variable "action" {}

resource "huaweicloud_dds_backup_stop" "test"{
  backup_id = var.backup_id
  action    = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `backup_id` - (Required, String, NonUpdatable) Specifies the ID of the backup.

* `action` - (Required, String, NonUpdatable) Specifies the action.
  The value only can be **stop**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
