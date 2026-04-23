---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topic_unsubscription"
description: |-
  Manages a resource to unsubscribe SMN topic within HuaweiCloud.
---

# huaweicloud_smn_topic_unsubscription

Manages a resource to unsubscribe SMN topic within HuaweiCloud.

## Example Usage

```hcl
variable "subscription_urn" {}

resource "huaweicloud_smn_topic_unsubscription" "test" {
  subscription_urn = var.subscription_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the SMN topic subscriber resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `subscription_urn` - (Required, String, NonUpdatable) Specifies the unique resource identifier of a subscriber.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `message` - The response information.
