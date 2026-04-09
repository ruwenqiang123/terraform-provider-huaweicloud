package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAllL7Rules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_elb_all_l7rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Beforce running test, prepare a loadbalancer with a listener, and the listener associate forword rule.
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAllL7Rules_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.compare_type"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.provisioning_status"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.invert"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.conditions.#"),

					resource.TestCheckOutput("rule_id_filter_useful", "true"),
					resource.TestCheckOutput("type_filter_useful", "true"),
					resource.TestCheckOutput("compare_type_filter_useful", "true"),
					resource.TestCheckOutput("value_filter_useful", "true"),
					resource.TestCheckOutput("provisioning_status_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAllL7Rules_basic = `
data "huaweicloud_elb_all_l7rules" "test" {}

locals {
  rule_id             = data.huaweicloud_elb_all_l7rules.test.rules[0].id
  type                = data.huaweicloud_elb_all_l7rules.test.rules[0].type
  compare_type        = data.huaweicloud_elb_all_l7rules.test.rules[0].compare_type
  value               = data.huaweicloud_elb_all_l7rules.test.rules[0].value
  provisioning_status = data.huaweicloud_elb_all_l7rules.test.rules[0].provisioning_status
}

data "huaweicloud_elb_all_l7rules" "rule_id_filter" {
  rule_id = [local.rule_id]
}

output "rule_id_filter_useful" {
  value = length(data.huaweicloud_elb_all_l7rules.rule_id_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_elb_all_l7rules.rule_id_filter.rules[*].id : v == local.rule_id]
  )
}

data "huaweicloud_elb_all_l7rules" "type_filter" {
  type = [local.type]
}

output "type_filter_useful" {
  value = length(data.huaweicloud_elb_all_l7rules.type_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_elb_all_l7rules.type_filter.rules[*].type : v == local.type]
  )
}

data "huaweicloud_elb_all_l7rules" "compare_type_filter" {
  compare_type = [local.compare_type]
}

output "compare_type_filter_useful" {
  value = length(data.huaweicloud_elb_all_l7rules.compare_type_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_elb_all_l7rules.compare_type_filter.rules[*].compare_type : v == local.compare_type]
  )
}

data "huaweicloud_elb_all_l7rules" "value_filter" {
  value = [local.value]
}

output "value_filter_useful" {
  value = length(data.huaweicloud_elb_all_l7rules.value_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_elb_all_l7rules.value_filter.rules[*].value : v == local.value]
  )
}

data "huaweicloud_elb_all_l7rules" "provisioning_status_filter" {
  provisioning_status = [local.provisioning_status]
}

output "provisioning_status_filter_useful" {
  value = length(data.huaweicloud_elb_all_l7rules.provisioning_status_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_elb_all_l7rules.provisioning_status_filter.rules[*].provisioning_status : v == local.provisioning_status]
  )
}
`
