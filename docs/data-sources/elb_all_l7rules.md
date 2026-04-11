---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_all_l7rules"
description: |-
  Use this data source to query the all forwarding rules of a specified region within HuaweiCloud.
---

# huaweicloud_elb_all_l7rules

Use this data source to query the all forwarding rules of a specified region within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_elb_all_l7rules" "test" {}
```

### Filter by rule ID

```hcl
variable "rule_id" {
  type = list(string)
}

data "huaweicloud_elb_all_l7rules" "test" {
  rule_id = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the load balancer topology.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, List) Specifies the forwarding rule ID list.

* `type` - (Optional, List) Specifies the match type.
  The value can be **HOST_NAME** or **PATH**.

* `compare_type` - (Optional, List) Specifies the match method.
  The valid values are as follows:
  + **EQUAL_TO**: Exact match.
  + **REGEX**: Regular expression match.
  + **STARTS_WITH**: Prefix match.

* `provisioning_status` - (Optional, List) Specifies the provisioning status of the forwarding rule.
  The value can only be **ACTIVE**.

* `value` - (Optional, List) Specifies the value of the match content.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  Default query the resource under all enterprise projects.
  If you need to query data for a specific enterprise project, you need set the corresponding to enterprise project ID.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of forwarding rules.
  The [rules](#l7rules_struct) structure is documented below.

<a name="l7rules_struct"></a>
The `rules` block supports:

* `id` - The forwarding rule ID.

* `type` - The forwarding rule type.

* `compare_type` - The forwarding rule matching method.

* `key` - The match content key.

* `value` - The match content value.

* `provisioning_status` - The provisioning status.

* `invert` - Whether reverse matching is supported.

* `conditions` - The name of the load balancer.
  The [conditions](#l7policies_l7rules_conditions) structure is documented below.

* `project_id` - The project ID.

* `created_at` - The creation time.

* `updated_at` - The update time.

<a name="l7policies_l7rules_conditions"></a>
The `conditions` block supports:

* `key` - The key of match item.

* `value` - The value of match item.
