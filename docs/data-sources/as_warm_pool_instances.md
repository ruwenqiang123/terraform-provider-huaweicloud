---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_warm_pool_instances"
description: |-
  Use this data source to get the instances in the warm pool of an AS group within HuaweiCloud.
---

# huaweicloud_as_warm_pool_instances

Use this data source to get the instances in the warm pool of an AS group within HuaweiCloud.

## Example Usage

```hcl
variable "scaling_group_id" {}

data "huaweicloud_as_warm_pool_instances" "test" {
  scaling_group_id = var.scaling_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the scaling policies.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the scaling group ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `warm_pool_instances` - The instances in the warm pool.
  The [warm_pool_instances](#as_warm_pool_instances) structure is documented below.

<a name="as_warm_pool_instances"></a>
The `warm_pool_instances` block supports:

* `id` - The instance ID.

* `instance_id` - The VM ID.

* `name` - The instance name.

* `status` - The instance status.
