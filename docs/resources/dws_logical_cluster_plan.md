---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_logical_cluster_plan"
description: |-
  Use this resource to manage a logical cluster scaling plan within HuaweiCloud.
---

# huaweicloud_dws_logical_cluster_plan

Use this resource to manage a logical cluster scaling plan within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "logical_cluster_plan_actions" {
  type = list(object({
    type     = string
    strategy = string
  }))
}

resource "huaweicloud_dws_logical_cluster_plan" "test" {
  cluster_id = var.cluster_id
  plan_type  = "once"

  dynamic "actions" {
    for_each = var.logical_cluster_plan_actions
    
    content {
      type     = actions.value.type
      strategy = actions.value.strategy
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the logical cluster plan is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to which
  the logical cluster plan belongs.

* `plan_type` - (Required, String, NonUpdatable) Specifies the plan type.  
  The valid values are as follows:
  + **once**
  + **periodicity**

* `actions` - (Required, List) Specifies the list of logical cluster plan actions.  
  The [actions](#dws_logical_cluster_plan_actions) object structure is documented below.

* `logical_cluster_name` - (Optional, String, NonUpdatable) Specifies the logical cluster name.  
  Conflict with Parameter `user`. They cannot be used together.

* `user` - (Optional, String, NonUpdatable) Specifies the user bound to the logical cluster.  
  Conflict with Parameter `logical_cluster_name`. They cannot be used together.

* `node_num` - (Optional, Int, NonUpdatable) Specifies the number of logical cluster nodes.

* `main_logical_cluster` - (Optional, String, NonUpdatable) Specifies the main logical cluster bound to
  the logical cluster.

* `start_time` - (Optional, String, NonUpdatable) Specifies the start time of the logical cluster plan,
  in Unix timestamp format.

* `end_time` - (Optional, String, NonUpdatable) Specifies the end time of the logical cluster plan,
  in Unix timestamp format.

* `enabled` - (Optional, Bool) Specifies whether the logical cluster plan is enabled.

<a name="dws_logical_cluster_plan_actions"></a>
The `actions` block supports:

* `type` - (Required, String) Specifies the action type.

* `strategy` - (Required, String) Specifies the strategy expression or timestamp for the action.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the logical cluster plan.

* `actions` - The list of logical cluster plan actions.  
  The [actions](#dws_logical_cluster_plan_actions_attr) object structure is documented below.

<a name="dws_logical_cluster_plan_actions_attr"></a>
The `actions` block supports:

* `id` - The ID of the action.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is `30` minutes.
* `update` - Default is `30` minutes.
* `delete` - Default is `30` minutes.

## Import

The DWS logical cluster plan resource can be imported using the `cluster_id` and `plan_id`, separated by a slash, e.g.

```bash
terraform import huaweicloud_dws_logical_cluster_plan.test <cluster_id>/<plan_id>
```
