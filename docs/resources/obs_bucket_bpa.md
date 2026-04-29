---
subcategory: "Object Storage Service (OBS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_obs_bucket_bpa"
description: |-
  Provides an OBS bucket public access block resource.
---

# huaweicloud_obs_bucket_bpa

Provides an OBS bucket public access block resource.

-> 1. Currently, the BPA capability of this bucket is only supported by some regions. Using this resource in regions that
  do not support it may result in errors.<br/>2. Deleting this resource will set all BPA configurations to **false**.

## Example Usage

```hcl
resource "huaweicloud_obs_bucket_bpa" "example" {
  bucket                  = "your_bucket_name"
  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `bucket` - (Required, String, NonUpdatable) Specifies the name of the bucket to set public access block.

* `block_public_acls` - (Optional, Bool) Specifies whether to lock public ACLs. When set to true, uploading objects with
  public ACLs is prohibited, and ACL modification APIs that set public ACLs are forbidden.
  Defaults to **false**.

* `ignore_public_acls` - (Optional, Bool) Specifies whether to ignore public ACLs. When set to true, public ACLs do not
  take effect during all OBS OpenAPI permission checks.
  Defaults to **false**.

* `block_public_policy` - (Optional, Bool) Specifies whether to lock public policies. When set to true, bucket policy
  modification APIs that set public policies are forbidden.
  Defaults to **false**.

* `restrict_public_buckets` - (Optional, Bool) Specifies whether to restrict account access. When set to true, during
  all OBS OpenAPI permission checks, if the bucket policy status is public, only cloud service accounts and the current
  account are allowed to access.
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the bucket name.

## Import

OBS bucket public access block can be imported using the bucket name, e.g.

```bash
$ terraform import huaweicloud_obs_bucket_bpa.test <bucket>
```
