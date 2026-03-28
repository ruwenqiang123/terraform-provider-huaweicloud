---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_factory_job_export"
description: |-
  Use this resource to export a DataArts Factory job to a specified OBS storage path within HuaweiCloud.
---

# huaweicloud_dataarts_factory_job_export

Use this resource to export a DataArts Factory job to a specified OBS storage path within HuaweiCloud.

-> This resource is a one-time action resource used to export a job package to a specified OBS storage path. Deleting
   this resource will not clear the export task, but will only remove the resource information from the tfstate file.

-> After the resource is created, a job export folder will be generated in the target OBS storage path, which contains
   the job export object (ZIP package, names `{job_name}.zip`). Generating the ZIP package takes some time (it will not
   be generated instantly after the resource is created).

## Example Usage

```hcl
variable "workspace_id" {}
variable "job_name" {}
variable "obs_storage_path" {}

resource "huaweicloud_dataarts_factory_job_export" "test" {
  workspace_id  = var.workspace_id
  job_name      = var.job_name
  export_depend = true
  export_status = "SUBMIT"
  obs_path      = var.obs_storage_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_name` - (Required, String, NonUpdatable) Specifies the name of the job to be exported.

* `obs_path` - (Required, String, NonUpdatable) Specifies the OBS path for storing the exported job package.
  For example, `obs://test_bucket/job_nodes/`.

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the workspace ID to which the job belongs.  
  If omitted, the default workspace (`default`) will be used.

* `export_depend` - (Optional, Bool, NonUpdatable) Specifies whether to export scripts and resources depended on
  by the job. If omitted, the default value defined by the API is used.

* `export_status` - (Optional, String, NonUpdatable) Specifies the job status to be exported.  
  The valid values are as follows:
  + **DEVELOP**
  + **SUBMIT**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `folder_path` - The folder path suffix parsed from the response `obsPath` relative to the input `obs_path`.  
  For example, if `obs_path` is `obs://dataarts-test/jobs-storage/` and the response `obsPath` is
  `obs://dataarts-test/jobs-storage/job_6278_exported/1774334941255`, this attribute will be
  `job_6278_exported/1774334941255`.
