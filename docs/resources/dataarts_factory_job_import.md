---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_job_import"
description: |-
  Use this resource to import a DataArts Factory job from an specified OBS storage path within HuaweiCloud.
---

# huaweicloud_dataarts_factory_job_import

Use this resource to import a DataArts Factory job from an OBS storage path within HuaweiCloud.

-> This resource is a one-time action resource used to import jobs from a specified OBS storage path. Deleting this
   resource will not clear the import task, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "obs_storage_path" {}

resource "huaweicloud_dataarts_factory_job_import" "test" {
  workspace_id     = var.workspace_id
  path             = var.obs_storage_path
  same_name_policy = "OVERWRITE"
  target_status    = "SAVED"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the job is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `path` - (Required, String, NonUpdatable) Specifies the OBS path where the job package is stored.  
  For example, `obs://myBucket/jobs.zip`.

* `same_name_policy` - (Optional, String, NonUpdatable) Specifies the duplicate name handling policy.  
  The valid values are as follows:
  + **SKIP**
  + **OVERWRITE**

* `target_status` - (Optional, String, NonUpdatable) Specifies the target status of imported jobs.  
  The valid values are as follows:
  + **SAVED**
  + **SUBMITTED**
  + **PRODUCTION**

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the workspace ID to which the imported jobs belong.  
  If omitted, the default workspace (`default`) will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
