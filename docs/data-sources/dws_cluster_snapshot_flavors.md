---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_snapshot_flavors"
description: |-
  Use this data source to query the DWS cluster flavors by snapshot ID within HuaweiCloud.
---

# huaweicloud_dws_cluster_snapshot_flavors

Use this data source to query the DWS cluster flavors by snapshot ID within HuaweiCloud.

## Example Usage

```hcl
variable "snapshot_id" {}
data "huaweicloud_dws_cluster_snapshot_flavors" "test" {
  snapshot_id = var.snapshot_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the cluster flavors are located.  
  If omitted, the provider-level region will be used.

* `snapshot_id` - (Required, String) Specifies the ID of the snapshot to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `flavors` - The list of cluster flavors for the snapshot.  
  The [flavors](#dws_snapshot_flavors) structure is documented below.

<a name="dws_snapshot_flavors"></a>
The `flavors` block supports:

* `id` - The flavor ID.

* `code` - The flavor code.

* `classify` - The flavor type.

* `scenario` - The flavor scenario.

* `version` - The flavor version.

* `status` - The flavor status.

* `default_capacity` - The default capacity of the flavor.

* `duplicate` - The number of replicas used by the flavor.

* `default_node` - The default number of nodes.

* `min_node` - The minimum number of nodes.

* `max_node` - The maximum number of nodes.

* `flavor_id` - The underlying flavor ID.

* `flavor_code` - The underlying flavor code.

* `volume_num` - The number of disks.

* `attribute` - The list of extended information.  
  The [attribute](#dws_snapshot_flavors_attribute) structure is documented below.

* `product_version_list` - The list of product versions supported by the flavor.  
  The [product_version_list](#dws_snapshot_flavors_product_version_list) structure is documented below.

* `volume_used` - The disk usage information of the snapshot source cluster.  
  The [volume_used](#dws_snapshot_flavors_volume_used) structure is documented below.

<a name="dws_snapshot_flavors_attribute"></a>
The `attribute` block supports:

* `code` - The extended information code.

* `value` - The extended information value.

<a name="dws_snapshot_flavors_product_version_list"></a>
The `product_version_list` block supports:

* `min_cn` - The minimum number of CN nodes supported by this version.

* `max_cn` - The maximum number of CN nodes supported by this version.

* `version_type` - The type of this version.

* `datastore_version` - The datastore version name.

<a name="dws_snapshot_flavors_volume_used"></a>
The `volume_used` block supports:

* `volume_type` - The disk type.

* `volume_num` - The number of disks.

* `capacity` - The available storage capacity of a single node, in GB.

* `volume_size` - The physical storage capacity of a single disk, in GB.
