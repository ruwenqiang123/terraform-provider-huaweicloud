---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_batch_delete_address_groups"
description: |-
  Manages a resource to batch delete address groups within HuaweiCloud.
---

# huaweicloud_cfw_batch_delete_address_groups

Manages a resource to batch delete address groups within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch delete address groups. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. After the successful execution of the resource, it is necessary to pay attention to the value of the `data`
  attribute. If the value of `data` is empty, it means that the current operation has not taken effect.

## Example Usage

```hcl
variable "object_id" {}
variable "set_ids" {
  type = list(string)
}

resource "huaweicloud_cfw_batch_delete_address_groups" "test" {
  object_id = var.object_id
  set_ids   = var.set_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID.

* `set_ids` - (Required, List, NonUpdatable) Specifies the IDs of address groups to be deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `object_id`.

* `data` - The address groups for batch deletion.

  The [data](#Batch_Delete_Address_Groups_Data) structure is documented below.

<a name="Batch_Delete_Address_Groups_Data"></a>
The `data` block supports:

* `name` - The name of the address group for batch deletion.

* `id` - The ID of the address group for batch deletion.
