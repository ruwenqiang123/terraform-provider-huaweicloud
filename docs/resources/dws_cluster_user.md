---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_user"
description: |-
  Manages a DWS cluster user or role resource within HuaweiCloud.
---

# huaweicloud_dws_cluster_user

Manages a DWS cluster user or role resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "cluster_id" {}
variable "user_name" {}
variable "user_password" {}

resource "huaweicloud_dws_cluster_user" "test" {
  cluster_id   = var.cluster_id
  type         = "user"
  name         = var.user_name
  password     = var.user_password
}
```

### Granting object permissions

```hcl
variable "cluster_id" {}
variable "user_name" {}
variable "user_password" {}
variable "user_type" {}

variable "grant_list" {
  type = list(object({
    type       = string
    database   = string
    privileges = list(object({
      permission = string
      grant_with = bool
    }))
  }))
}

resource "huaweicloud_dws_cluster_user" "with_grants" {
  cluster_id = var.cluster_id
  type       = var.user_type
  name       = var.user_name
  password   = var.user_password

  dynamic "grant_list" {
    for_each = var.grant_list
    content {
      type     = grant_list.value.type
      database = grant_list.value.database

      dynamic "privileges" {
        for_each = grant_list.value.privileges
        content {
          permission = privileges.value.permission
          grant_with = privileges.value.grant_with
        }
      }
    }
  }
}
```

### Configuring user attributes

```hcl
variable "cluster_id" {}
variable "user_name" {}
variable "user_password" {}
variable "user_login" {}
variable "user_create_role" {}
variable "user_create_db" {}
variable "user_inherit" {}
variable "user_conn_limit" {}
variable "user_lock" {}

resource "huaweicloud_dws_cluster_user" "attributes_only" {
  cluster_id  = var.cluster_id
  type        = "user"
  name        = var.user_name
  password    = var.user_password
  login       = var.user_login
  create_role = var.user_create_role
  create_db   = var.user_create_db
  inherit     = var.user_inherit
  conn_limit  = var.user_conn_limit
  lock        = var.user_lock
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the cluster user is located.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `cluster_id` - (Required, String, NonUpdatable) Specifies the ID of the DWS cluster.

* `name` - (Required, String, NonUpdatable) Specifies the name of the cluster user or role.

* `type` - (Optional, String) Specifies the object type.
  The valid values are as follows:
  + **user**
  + **role**

  Defaults to **user**.

* `password` - (Optional, String, NonUpdatable) Specifies the password of the user.
  Changing this parameter will create a new resource.

  -> Required if the `type` is **user**.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the user or role.

* `password_disable` - (Optional, Bool, NonUpdatable) Specifies whether to disable password authentication.

* `logical_cluster` - (Optional, String, NonUpdatable) Specifies the logical cluster name.

* `cascade` - (Optional, Bool) Specifies whether to cascade delete dependencies when deleting the user or role.

* `grant_list` - (Optional, List, NonUpdatable) Specifies the set of grants.
  The [grant_list](#dws_cluster_user_grant_list) structure is documented below.

* `login` - (Optional, Bool) Specifies whether to allow the user to log in.

* `create_role` - (Optional, Bool) Specifies whether to grant the permission to create roles.

* `create_db` - (Optional, Bool) Specifies whether to grant the permission to create databases.

* `system_admin` - (Optional, Bool) Specifies whether to grant the system administrator permission.

* `audit_admin` - (Optional, Bool) Specifies whether to grant the audit administrator permission.

* `inherit` - (Optional, Bool) Specifies whether to inherit permissions from roles.

* `use_ft` - (Optional, Bool) Specifies whether to grant the external table permission.

* `conn_limit` - (Optional, Int) Specifies the maximum number of concurrent connections.

* `replication` - (Optional, Bool) Specifies whether to grant the replication permission.

* `valid_begin` - (Optional, String) Specifies the valid begin time.

* `valid_until` - (Optional, String) Specifies the valid until time.

* `lock` - (Optional, Bool) Specifies whether the user is locked.

<a name="dws_cluster_user_grant_list"></a>
The `grant_list` block supports:

* `type` - (Required, String) The object type.
  The valid values are as follows:
  + **DATABASE**
  + **SCHEMA**
  + **TABLE**
  + **VIEW**
  + **COLUMN**
  + **FUNCTION**
  + **SEQUENCE**
  + **NODEGROUP**

* `database` - (Optional, String) The database name.

* `schema_name` - (Optional, String) The schema name.

* `object_name` - (Optional, String) The object name.

* `all_object` - (Optional, Bool) Whether all objects are effective.

* `future` - (Optional, Bool) Whether future objects are effective.

* `future_object_owners` - (Optional, String) The owners of future objects.

* `column_names` - (Optional, List) The set of column names.

* `privileges` - (Required, List) The set of privileges.
  The [privileges](#dws_cluster_user_privileges) structure is documented below.

<a name="dws_cluster_user_privileges"></a>
The `privileges` block supports:

* `permission` - (Required, String) The privilege name.
  The valid values are as follows:
  + When `type` is **DATABASE**:
    - **CREATE**
    - **CONNECT**
    - **TEMPORARY**
    - **TEMP ALL PRIVILEGES**
  + When `type` is **SCHEMA**:
    - **CREATE**
    - **USAGE**
    - **ALTER**
    - **DROP ALL PRIVILEGES**
  + When `type` is **TABLE**:
    - **SELECT**
    - **INSERT**
    - **UPDATE**
    - **DELETE**
    - **TRUNCATE**
    - **REFERENCES**
    - **TRIGGER**
    - **ANALYZE**
    - **VACUUM**
    - **ALTER**
    - **DROP ALL PRIVILEGES**
  + When `type` is **VIEW**:
    - **SELECT**
    - **INSERT**
    - **UPDATE**
    - **DELETE**
    - **TRUNCATE**
    - **REFERENCES**
    - **TRIGGER**
    - **ANALYZE**
    - **VACUUM**
    - **ALTER**
    - **DROP ALL PRIVILEGES**
  + When `type` is **COLUMN**:
    - **SELECT**
    - **INSERT**
    - **UPDATE**
    - **REFERENCES ALL PRIVILEGES**
  + When `type` is **FUNCTION**:
    - **EXECUTE ALL PRIVILEGES**
  + When `type` is **SEQUENCE**:
    - **SELECT**
    - **UPDATE**
    - **USAGE ALL PRIVILEGES**
  + When `type` is **NODEGROUP**:
    - **CREATE**
    - **USAGE**
    - **COMPUTE ALL PRIVILEGES**

* `grant_with` - (Required, Bool) Whether the grant option is included.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The DWS cluster user resource can be imported using the region where the cluster is located, the ID of the DWS cluster,
and the name of the cluster user or role, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_dws_cluster_user.test <region>/<cluster_id>/<name>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `password`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dws_cluster_user" "test" {
  ...

  lifecycle {
    ignore_changes = [
      type,
      password,
      cascade,
    ]
  }
}
```
