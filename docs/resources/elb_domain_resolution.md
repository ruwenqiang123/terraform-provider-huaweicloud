---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_domain_resolution"
description: |-
  Manages an ELB domain resolution resource within HuaweiCloud.
---

# huaweicloud_elb_domain_resolution

Manages an ELB domain resolution resource within HuaweiCloud.

## Example Usage

```hcl
variable "loadbalancer_id" {}
variable "public_dns_zone_name" {}
variable "private_dns_zone_name" {}

resource "huaweicloud_elb_domain_resolution" "test" {
  loadbalancer_id            = var.loadbalancer_id
  public_domain_name_enable  = true
  public_dns_zone_name       = var.public_dns_zone_name
  private_domain_name_enable = true
  private_dns_zone_name      = var.private_dns_zone_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `loadbalancer_id` - (Required, String, NonUpdatable) Specifies the ID of the load balancer.

* `public_domain_name_enable` - (Optional, Bool) Specifies whether to enable public domain name resolution.
  The valid values are as follows:
  + **true**: Enable public domain name resolution.
  + **false**: Disable public domain name resolution. (Default value).

* `private_domain_name_enable` - (Optional, Bool) Specifies whether to enable private domain name resolution.
  The valid values are as follows:
  + **true**: Enable private domain name resolution.
  + **false**: Disable private domain name resolution. (Default value).

-> Creating the resource, At least one of `public_domain_name_enable` and `private_domain_name_enable` needs to
   set to **true**.

* `public_dns_zone_name` - (Optional, String) Specifies the public zone that will be used to generate public
  domain names for the load balancer.

  -> 1. This parameter is valid and required only when `public_domain_name_enable` is set to **true**.
     <br/>2. The public zone you specified must exist on the DNS servier.
     <br/>3. The public domain names can only be generated based on public zones.

* `public_dns_record_set_ttl` - (Optional, Int) Specifies the cache duration of the public record set on a
  local DNS server，in seconds.
  The valid value ranges from `1` to `2,147,483,647`, Default to `300`.
  The longer the duration is, the slower the update takes effect. If your service address changes frequently,
  set TTL to a smaller value. Otherwise, set TTL to a larger value.

* `private_dns_zone_name` - (Optional, String) Specifies the zone that will be used to generate private
  domain names for the load balancer.

  -> 1. This parameter is valid and required only when `private_domain_name_enable` is set to **true**.
     <br/>2. The private zone you specified must exist on the DNS servier.
     <br/>3. Both public and private zones can be used to generate private domain names for the load balancer.
     The zone type is defined by `private_dns_zone_type`.

* `private_dns_zone_type` - (Optional, String) Specifies the type of the zone that will be used to generate
  private domain names for the load balancer.
  The valid values are as follows:
  + **private**: A private zone will be used.
  + **public**: A public zone will be used.

* `private_dns_record_set_ttl` - (Optional, Int) Specifies the cache duration of the private record set on a
  local DNS server, in seconds.
  The valid value ranges from `1` to `2,147,483,647`, Default to `300`.
  The longer the duration is, the slower the update takes effect. If your service address changes frequently,
  set TTL to a smaller value. Otherwise, set TTL to a larger value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the load balancer ID.

* `ips` - The domain name resolution information of a load balancer.
  The [ips](#domain_resolution_ips_struct) structure is documented below.

<a name="domain_resolution_ips_struct"></a>
The `ips` block supports:

* `enable` - Whether domain name resolution is enabled for the IP address.

* `ip_address` - The IP address.

* `type` - The IP address type.

* `domain_name` - The domain name associated with the IP address.

* `created_at` - The create time.

* `updated_at` - The update time.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_domain_resolution.test <id>
```
