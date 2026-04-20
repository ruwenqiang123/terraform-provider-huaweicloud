---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_internet_bandwidths_by_tags"
description: |-
  Use this data source to get the list of global internet bandwidths filtered by tags.
---

# huaweicloud_global_internet_bandwidths_by_tags

Use this data source to get the list of global internet bandwidths filtered by tags.

## Example Usage

```hcl
data "huaweicloud_global_internet_bandwidths_by_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `tags` - (Optional, List) Specifies the tags to filter global internet bandwidths.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the tag key. The key of the same resource cannot be duplicated.
  If a created predefined tag is exactly the same as an existing predefined tag, it will overwrite the existing one.
  The maximum length is `128` characters, which can consist of English letters, digits, underscores (_), hyphens (-),
  and Chinese characters.

* `value` - (Optional, String) Specifies the tag value. The maximum length of the value is `255` characters, which can
  consist of English letters, digits, underscores (_), dots (.), hyphens (-), and Chinese characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - The list of global internet bandwidths that match the specified tags.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `resource_id` - The ID of the global internet bandwidth resource.

* `resource_name` - The name of the global internet bandwidth resource.

* `resource_detail` - The detailed information of the global internet bandwidth resource.

* `tags` - The tags associated with the resource.

  The [tags](#resource_tags_struct) structure is documented below.

<a name="resource_tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.
