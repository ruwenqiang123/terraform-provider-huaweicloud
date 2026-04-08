---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_public_domain_associate"
description: |-
  Manages a public domain for the DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_public_domain_associate

Manages a public domain for the DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "public_domain_name" {}

resource "huaweicloud_dws_cluster_public_domain_associate" "test" {
  cluster_id  = var.cluster_id
  domain_name = var.public_domain_name
  ttl         = 1000
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cluster (to which the public domain belongs) is
  located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the cluster ID to which the public domain belongs.

* `domain_name` - (Required, String) Specifies the public domain name.

* `ttl` - (Optional, Int) Specifies the cache period of the SOA record set, in seconds.  
  The valid value is range from `300` to `2,147,483,647`, if omitted, the service side default value will be used.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID, also `cluster_id`.

## Import

The resource can be imported using the `id` (also `cluster_id`), e.g.

```bash
$ terraform import huaweicloud_dws_cluster_public_domain_associate.test <id>
```
