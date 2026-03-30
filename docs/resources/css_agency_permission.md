---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_agency_permission"
description: |-
  Manages an agency permission resource within HuaweiCloud.
---

# huaweicloud_css_agency_permission

Manages an agency permission resource within HuaweiCloud.

-> 1. Use this resource, need the built-in CSS agency exists,  the system will remove high-risk permissions and
   sets the minimal permissions. also is the administrator, fullaccess permission of the selected agency will be
   deleted and minimal permissions will be granted instead.
   <br/>2. This resource is a one-time action resource. Deleting this resource will not clear the corresponding
   request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "domain_id" {}
variable "domain_name" {}
variable "type" {}

resource "huaweicloud_css_agency_permission" "test" {
  domain_id   = var.domain_id
  domain_name = var.domain_name
  type        = var.type
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, String, NonUpdatable) Specifies the account ID.

* `domain_name` - (Required, String, NonUpdatable) Specifies the account name.

* `type` - (Required, String, NonUpdatable) Specifies the type of the agency.
  The valid values are as follows:
  + **obs**: Indicates agency permissions required for creating snapshots and log backups.
  + **vpc**: Indicates agency permissions required for version upgrade, AZ change, scale-in, and node replacement.
  + **elb**: Indicates agency permissions required for using the alerting plug-in.
  + **smn**: Indicates agency permissions required for using load balancing.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
