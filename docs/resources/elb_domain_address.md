---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_domain_address"
description: |-
  Manages an ELB domain IP address resource within HuaweiCloud.
---

# huaweicloud_elb_domain_address

Manages an ELB domain IP address resource within HuaweiCloud.

## Example Usage

```hcl
variable "loadbalancer_id" {}
variable "ip_address" {}

resource "huaweicloud_elb_domain_address" "test" {
  loadbalancer_id = var.loadbalancer_id
  ip_address      = var.ip_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `loadbalancer_id` - (Required, String, NonUpdatable) Specifies the ID of the load balancer.

* `ip_address` - (Required, String, NonUpdatable) Specifies the IP address.
  It can be an IPv4 or IPv6 address.

  -> The IP address must be a private or public IP address of the load balancer.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the load balancer ID and IP address.
  e.g. **1cab5d49-fb3d-4b99-8b05-c27d661152fd/192.168.1.105**

* `enable` - Whether domain resolution is enabled for the IP address.

* `type` - The IP address type.
  + **vip**: Indicates private IP address.
  + **eip**: Indicates public IP address.

* `domain_name` - The domain name associated with the IP address.

* `created_at` - The create time.

* `updated_at` - The update time.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_domain_address.test <id>
```
