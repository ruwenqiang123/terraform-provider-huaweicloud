package smn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmnTopicsAssociateResources_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_smn_topics_associate_resources.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Beforce running test, prepare a SMN topic and enabled the LTS.
			acceptance.TestAccPrecheckSmnFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnTopicsAssociateResources_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "topics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.topic_id"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.attributes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.logtanks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "topics.0.update_time"),

					resource.TestCheckOutput("topic_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceSmnTopicsAssociateResources_basic = `
data "huaweicloud_smn_topics_associate_resources" "test" {}

# Filter by topic_id
locals {
  topic_id = data.huaweicloud_smn_topics_associate_resources.test.topics[0].topic_id
}

data "huaweicloud_smn_topics_associate_resources" "filter_by_topic_id" {
  topic_id = local.topic_id
}

locals {
  topic_id_filter_result = [
    for v in data.huaweicloud_smn_topics_associate_resources.filter_by_topic_id.topics[*].topic_id : v == local.topic_id
  ]
}

output "topic_id_filter_is_useful" {
  value = alltrue(local.topic_id_filter_result) && length(local.topic_id_filter_result) > 0
}

# Filter by name
locals {
  topic_name = data.huaweicloud_smn_topics_associate_resources.test.topics[0].name
}

data "huaweicloud_smn_topics_associate_resources" "filter_by_name" {
  name = local.topic_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_smn_topics_associate_resources.filter_by_name.topics[*].name : v == local.topic_name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Empty result verification
data "huaweicloud_smn_topics_associate_resources" "non_exist" {
  name = "non-exist-topic-name"
}

output "non_exist_is_zero" {
  value = length(data.huaweicloud_smn_topics_associate_resources.non_exist.topics) == 0
}
`
