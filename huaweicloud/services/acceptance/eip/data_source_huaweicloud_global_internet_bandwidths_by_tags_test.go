package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalInternetBandwidthsByTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_global_internet_bandwidths_by_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalInternetBandwidthsByTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.key", "foo"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.value", "bar"),
					resource.TestCheckOutput("is_tag_filter_useful", "true"),
					resource.TestCheckOutput("is_key_only_filter_useful", "true"),
					resource.TestCheckOutput("is_no_tag_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGlobalInternetBandwidthsByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_internet_bandwidths_by_tags" "test" {
  depends_on = [huaweicloud_global_internet_bandwidth.test]

  tags {
    key   = "foo"
    value = "bar"
  }
}

# Filter by multiple tags (should return empty)
data "huaweicloud_global_internet_bandwidths_by_tags" "multiple_tags_filter" {
  depends_on = [huaweicloud_global_internet_bandwidth.test]

  tags {
    key   = "foo"
    value = "bar"
  }

  tags {
    key   = "non-exist-key"
    value = "non-exist-value"
  }
}

output "is_tag_filter_useful" {
  value = length(data.huaweicloud_global_internet_bandwidths_by_tags.multiple_tags_filter.resources) == 0
}

# Filter by tag key only (without value)
data "huaweicloud_global_internet_bandwidths_by_tags" "key_only_filter" {
  depends_on = [huaweicloud_global_internet_bandwidth.test]

  tags {
    key = "foo"
  }
}

output "is_key_only_filter_useful" {
  value = length(data.huaweicloud_global_internet_bandwidths_by_tags.key_only_filter.resources) > 0
}

# Query without tags (should return all resources)
data "huaweicloud_global_internet_bandwidths_by_tags" "no_tag_filter" {
  depends_on = [huaweicloud_global_internet_bandwidth.test]
}

output "is_no_tag_filter_useful" {
  value = length(data.huaweicloud_global_internet_bandwidths_by_tags.no_tag_filter.resources) > 0
}
`, testAccInternetBandwidth_basic(name))
}
