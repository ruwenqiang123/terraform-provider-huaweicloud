package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLoadbalancerTopology_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_elb_loadbalancer_topology.test"
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
				Config: testAccDataSourceLoadbalancerTopology_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "edges.#"),
					resource.TestCheckResourceAttrSet(dataSource, "labels.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.loadbalancers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.loadbalancers.0.guaranteed"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.loadbalancers.0.availability_zone_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.listeners.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.listeners.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.listeners.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.listeners.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.pools.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.pools.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.pools.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "edges.0.source"),
					resource.TestCheckResourceAttrSet(dataSource, "edges.0.target"),
					resource.TestCheckResourceAttrSet(dataSource, "edges.0.source_type"),
					resource.TestCheckResourceAttrSet(dataSource, "edges.0.target_type"),
					resource.TestCheckResourceAttrSet(dataSource, "labels.0.l7policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "labels.0.l7policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "labels.0.l7policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "labels.0.l7policies.0.priority"),

					resource.TestCheckOutput("listener_id_filter_useful", "true"),
					resource.TestCheckOutput("listener_name_filter_useful", "true"),
					resource.TestCheckOutput("listener_protocol_filter_useful", "true"),
					resource.TestCheckOutput("pool_id_filter_useful", "true"),
					resource.TestCheckOutput("pool_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceLoadbalancerTopology_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_elb_loadbalancer_topology" "test" {
  loadbalancer_id = "%[1]s"
}

locals {
  listener_id       = data.huaweicloud_elb_loadbalancer_topology.test.nodes[0].listeners[0].id
  listener_name     = data.huaweicloud_elb_loadbalancer_topology.test.nodes[0].listeners[0].name
  listener_protocol = data.huaweicloud_elb_loadbalancer_topology.test.nodes[0].listeners[0].protocol
  pool_id           = data.huaweicloud_elb_loadbalancer_topology.test.nodes[0].pools[0].id
  pool_name         = data.huaweicloud_elb_loadbalancer_topology.test.nodes[0].pools[0].name

}

data "huaweicloud_elb_loadbalancer_topology" "listener_id_filter" {
  loadbalancer_id = "%[1]s"	
  listener_id     = local.listener_id
}

output "listener_id_filter_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_topology.listener_id_filter.nodes[0].listeners) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_topology.listener_id_filter.nodes[0].listeners[*].id : v == local.listener_id]
  )
}

data "huaweicloud_elb_loadbalancer_topology" "listener_name_filter" {
  loadbalancer_id = "%[1]s"	
  listener_name   = local.listener_name
}

output "listener_name_filter_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_topology.listener_name_filter.nodes[0].listeners) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_topology.listener_name_filter.nodes[0].listeners[*].name : v == local.listener_name]
  )
}

data "huaweicloud_elb_loadbalancer_topology" "listener_protocol_filter" {
  loadbalancer_id   = "%[1]s"	
  listener_protocol = local.listener_protocol
}

output "listener_protocol_filter_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_topology.listener_protocol_filter.nodes[0].listeners) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_topology.listener_protocol_filter.nodes[0].listeners[*].protocol : v == local.listener_protocol]
  )
}

data "huaweicloud_elb_loadbalancer_topology" "pool_id_filter" {
  loadbalancer_id = "%[1]s"	
  pool_id         = local.pool_id
}

output "pool_id_filter_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_topology.pool_id_filter.nodes[0].pools) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_topology.pool_id_filter.nodes[0].pools[*].id : v == local.pool_id]
  )
}

data "huaweicloud_elb_loadbalancer_topology" "pool_name_filter" {
  loadbalancer_id = "%[1]s"	
  pool_name       = local.pool_name
}

output "pool_name_filter_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_topology.pool_name_filter.nodes[0].pools) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_topology.pool_name_filter.nodes[0].pools[*].name : v == local.pool_name]
  )
}
`, acceptance.HW_ELB_LOADBALANCER_ID)
}
