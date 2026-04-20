---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topic_subscription"
description: |-
  Manages a resource to subscribe SMN topic within HuaweiCloud.
---

# huaweicloud_smn_topic_subscription

Manages a resource to subscribe SMN topic within HuaweiCloud.

## Example Usage

```hcl
variable "token" {}
variable "topic_urn" {}

resource "huaweicloud_smn_topic_subscription" "test" {
  token     = var.token
  topic_urn = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SMN topic subscriber resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `token` - (Required, String, NonUpdatable) Specifies the token information of a subscribed topic.

* `topic_urn` - (Optional, String, NonUpdatable) Specifies the unique resource identifier of the topic.

* `endpoint` - (Optional, String, NonUpdatable) Specifies the address of the subscription endpoint.

-> Exactly one of `topic_urn` or `endpoint` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `subscription_urn` - The unique resource identifier of a subscriber.
