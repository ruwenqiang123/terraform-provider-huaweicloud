---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_exception_rule"
description: |-
  Manages a exception rule under the DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_exception_rule

Manages a exception rule under the DWS cluster within HuaweiCloud.

~> Optional configurations removed during the update process will have their values ​​reset to the corresponding
   unrestricted value (-1).

## Example Usage

```hcl
variable "cluster_id" {}
variable "rule_name" {}
variable "exception_rule_configurations" {
  description = "The configurations of the exception rule."
  type        = list(object({
    key   = string
    value = string
  }))

  default = [
    {
      key   = "action"
      value = "penalty"
    },
    {
      key   = "spillsize"
      value = "300"
    },
    {
      key   = "bandwidth"
      value = "400"
    },
    {
      key   = "memsize"
      value = "500"
    },
    {
      key   = "broadcastsize"
      value = "600"
    },
  ]
}

resource "huaweicloud_dws_cluster_exception_rule" "test" {
  cluster_id = var.cluster_id
  name       = var.rule_name

  dynamic "configurations" {
    for_each = var.exception_rule_configurations

    content {
      key   = configurations.value.key
      value = configurations.value.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the exception rule is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the cluster to which the exception rule belongs.

* `name` - (Required, String, NonUpdatable) Specifies the name of the exception rule.  
  The valid value is limited from `3` to `28` characters, and only lowercase letters, digits, underscores (_) are
  allowed.

* `configurations` - (Required, List) Specifies the list of exception rule configurations.
  The [configurations](#dws_cluster_exception_rule_configurations) structure is documented below.

<a name="dws_cluster_exception_rule_configurations"></a>
The `configurations` block supports:

* `key` - (Required, String) Specifies the key of the exception rule configuration.

* `value` - (Required, String) Specifies the value of the exception rule configuration.  
  The valid keys and corresponding values are as follows:
  + **action** (Required configuration key): Operation upon exception.  
    The valid values are **abort** and **penalty**.
  + **blocktime**: Job blocking duration in seconds, which isthe total time consumed by the global andlocal concurrent
    queuing duration. The valueis an integer ranging from `1` to `2,147,483,647`. The value `-1` indicates nolimit.
  + **elapsedtime**: Job execution time in seconds. The value is an integer ranging from `1` to `2,147,483,647`.
    The value `-1` indicates no limit.
  + **allcputime**: Total CPU time spent in executing jobs on all DNs, in seconds.
    The value is an integer ranging from `1` to `2,147,483,647`. The value `-1` indicates no limit.
  + **cpuskewpercent**: Total CPU time skew rate on all DNs, in percentage (%).  
    The calculation formula is as follows:  
    **(Maximum CPU time of a statement across all DNs) × 100 / (Maximum CPU time of a statement across all DNs)**.  
    To set this rule, you need to set Interval for Checking CPU Skew Rate. The value is an integer in the range `1` to
    `100`. The value `-1` indicates no limit.
  + **cpuavgpercent**: Average CPU usage per DN, in percentage (%).  
    The calculation formula is as follows:  
    <!-- markdownlint-disable MD013 -->
    **(Total CPU time of a job on all DNs × 100) / (Job execution duration × Maximum CPUs on all DNs available to the resource pool associated with the job)**
    <!-- markdownlint-enable MD013 -->
    The value range is `1` to `100`. The value `-1` indicates no limit.
  + **spillsize**: Maximum allowed data (in MB) spilled to disk on a DN. The value range is `1` to `2,147,483,647`.  
    The value `-1` indicates no limit.
  + **bandwidth**: Maximum network bandwidth (in MB) that can be used by a job on a single DN.
    The value range is `1` to `2,147,483,647`. The value `-1` indicates no limit.
  + **memsize**: Maximum memory size (in MB) that can be used by a job on a single DN.
    The value range is `1` to `2,147,483,647`. The value `-1` indicates no limit.
  + **broadcastsize**: Maximum amount of large table broadcast data of a job on a single DN, in MB.
    The value is an integer ranging from `1` to `2,147,483,647`. The value `-1` indicates no limit.

  -> Except for **action** key, at least one of the other configurations must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in the format of `<cluster_id>/<name>`.

## Import

The resource can be imported using the `cluster_id` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dws_cluster_exception_rule.test <cluster_id>/<name>
```
