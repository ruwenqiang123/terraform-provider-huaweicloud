---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topics_associate_resources"
description: |-
  Use this data source to query the SMN topics with associated resources under a specified region within HuaweiCloud.
---

# huaweicloud_smn_topics_associate_resources

Use this data source to query the SMN topics with associated resources under a specified region within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_smn_topics_associate_resources" "test" {}
```

### Filter by topic name

```hcl
variable "topic_name" {}

data "huaweicloud_smn_topics_associate_resources" "test" {
  name = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resources.
  If omitted, the provider-level region will be used.

* `topic_id` - (Optional, String) Specifies the topic ID for exact matching.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `name` - (Optional, String) Specifies the topic name for exact matching.

* `fuzzy_name` - (Optional, String) Specifies the topic name for fuzzy matching.

* `fuzzy_display_name` - (Optional, String) Specifies the topic display name for fuzzy matching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics` - The list of topics with associated resources.  
  The [topics](#topics_struct) structure is documented below.

<a name="topics_struct"></a>
The `topics` block supports:

* `topic_urn` - The topic URN.

* `name` - The name of the topic.

* `display_name` - The display name of the topic.

* `push_policy` - The message push policy.  
  The valid values are as follows:
  + **0**: Send failed, keep to failed queue.
  + **1**: Directly discard failed messages.

* `enterprise_project_id` - The enterprise project ID.

* `topic_id` - The topic ID.

* `create_time` - The creation time.
  The time format is UTC time, **YYYY-MM-DDTHH:MM:SSZ**.

* `update_time` - The update time.
  The time format is UTC time, **YYYY-MM-DDTHH:MM:SSZ**.

* `tags` - The list of tags.  
  The [tags](#topics_tags_struct) structure is documented below.

* `attributes` - The topic access policy attributes.  
  The [attributes](#topics_attributes_struct) structure is documented below.

* `logtanks` - The cloud log information list.  
  The [logtanks](#topics_logtanks_struct) structure is documented below.

<a name="topics_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.

<a name="topics_attributes_struct"></a>
The `attributes` block supports:

* `access_policy` - The access policy of the topic.

* `create_time` - The creation time of the access policy.
  The time format is UTC time, **YYYY-MM-DDTHH:MM:SSZ**.

* `update_time` - The update time of the access policy.
  The time format is UTC time, **YYYY-MM-DDTHH:MM:SSZ**.

<a name="topics_logtanks_struct"></a>
The `logtanks` block supports:

* `id` - The ID of the cloud log.

* `log_group_id` - The ID of the log group.

* `log_stream_id` - The ID of the log stream.

* `create_time` - The creation time.
  The time format is UTC time, **YYYY-MM-DDTHH:MM:SSZ**.

* `update_time` - The update time.
  The time format is UTC time, **YYYY-MM-DDTHH:MM:SSZ**.
