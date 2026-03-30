---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_database_schema_adjust_action"
description: |-
  Use this resource to update the space limit of a database schema in a DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_database_schema_adjust_action

Use this resource to update the space limit of a database schema in a DWS cluster within HuaweiCloud.

-> This resource is only a one-time action resource for updating the schema space limit. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "cluster_id" {}
variable "database_name" {}
variable "schema_name" {}
variable "permission_space" {}

resource "huaweicloud_dws_database_schema_adjust_action" "test" {
  cluster_id = var.cluster_id
  database   = var.database_name
  schema     = var.schema_name
  perm_space = var.permission_space
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DWS cluster is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the DWS cluster.

* `database` - (Required, String, NonUpdatable) Specifies the name of the database.

* `schema` - (Required, String, NonUpdatable) Specifies the name of the schema.  
  In PostgreSQL-compatible databases (such as DWS), a schema is a namespace within a database
  that holds objects (for example, tables and views).  
  For instance, in the default setup you often see the `public` schema under a database.

* `perm_space` - (Required, Int, NonUpdatable) Specifies the space threshold of the schema.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the action resource.
