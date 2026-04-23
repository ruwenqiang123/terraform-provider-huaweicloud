---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_private_module_version"
description: |-
  Manages a RFS private module version resource within HuaweiCloud.
---

# huaweicloud_rfs_private_module_version

Manages a RFS private module version resource within HuaweiCloud.

## Example Usage

```hcl
variable "module_name" {}
variable "module_version" {}
variable "module_uri" {}

resource "huaweicloud_rfs_private_module_version" "test" {
  module_name    = var.module_name
  module_version = var.module_version
  module_uri     = var.module_uri
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource is located.
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `module_name` - (Required, String, NonUpdatable) Specifies the name of a private module. The name is unique within its
  domain and region. Only letters, digits, underscores (_), and hyphens (-) are allowed. It is case-sensitive and must
  start with a letter.

* `module_version` - (Required, String, NonUpdatable) Specifies the module version number. Each version number can be
  defined by following the semantic versioning format.

* `module_uri` - (Required, String, NonUpdatable) Specifies the address for accessing the module package in OBS. Modules
  allow you to bundle reusable code together. The OBS address can be accessed across regions of the same type. Regions
  are classified into universal regions and dedicated regions. A universal region provides universal cloud services for
  common tenants. A dedicated region provides specific services for specific tenants. The module package must be
  ZIP-formatted. To ensure the validity of the module package, the module package:
  + Cannot contain files whose names end with ".tfvars".
  + Cannot exceed `1` MB in size before and after decompression.
  + Cannot contain more than `100` files.
  + Cannot contain file paths starting with a forward slash (/).
  + Disallow spaces (" "), single periods ("."), or double periods ("..") between separators in file paths.
  + Allow each file path to be up to `2,048` characters long.
  + Allow each file name to be up to `255` characters long.
  + Must include at least one template file (whose name ends with ".tf" or ".tf.json").

  -> Modules do not support encryption for sensitive data. The module package associated with the module_uri is used,
  logged, displayed, and stored in plaintext by RFS.

* `module_id` - (Optional, String, NonUpdatable) Specifies the ID of a private module.

* `version_description` - (Optional, String, NonUpdatable) Specifies the description of a module version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The format is `<module_name>/<module_version>`.

* `create_time` - The creation time of a private module version. It is represented in UTC format
  (YYYY-MM-DDTHH:mm:ss.SSSZ), such as **1970-01-01T00:00:00.000Z**.

## Import

Private module version can be imported using `module_name` and `module_version` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_rfs_private_module_version.test <module_name>/<module_version>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `module_uri`. It is generally recommended running `terraform plan` after
importing a resource. You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_rfs_private_module_version" "test" {
  ...

  lifecycle {
    ignore_changes = [
      module_uri,
    ]
  }
}
```
