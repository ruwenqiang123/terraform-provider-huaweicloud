---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_backup_configurations"
description: |-
  Use this data source to query the extended backup configurations.
---

# huaweicloud_gaussdb_backup_configurations

Use this data source to query the extended backup configurations.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_backup_configurations" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rate_limit` - The backup flow control.

* `file_split_size` - The shard size.

* `prefetch_block` - The number of differential prefetch pages.

* `enable_standby_backup` - Whether to enable standby node backup.

* `close_compression` - Whether to disable backup compression.

* `default_backup_method` - The default backup method.
  + **PHYSICAL_BACKUP**: Indicates physical backup.
  + **EBACKUP**: Indicates snapshot backup.

* `default_backup_media_type` - The default backup storage medium.

* `backup_parallel_degree` - The number of concurrent backup tasks.
