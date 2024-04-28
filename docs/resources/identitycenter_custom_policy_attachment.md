---
subcategory: "IAM Identity Center"
---

# huaweicloud_identitycenter_custom_policy_attachment

Manages an Identity Center custom policy attachment resource within HuaweiCloud.

-> **NOTE:** Only one custom policy can be attached for a permission set, and it will be covered if another custom
strategy is attached.

## Example Usage

```hcl
variable "permission_set_id" {}

data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_custom_policy_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = var.permission_set_id
  custom_policy     = jsonencode(
    {
      "Version":"5.0",
      "Statement":[
        {
          "Effect":"Allow",
          "Action":[
            "eps:resources:add"
          ]
        }
      ]
    }
  )
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center instance.

  Changing this parameter will create a new resource.

* `permission_set_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center permission set.

  Changing this parameter will create a new resource.

* `custom_policy` - (Required, String) Specifies the custom policy to attach to a permission set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The Identity Center custom policy attachment can be imported using the `instance_id` and `permission_set_id` separated
by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_custom_policy_attachment.test <instance_id>/<permission_set_id>
```
