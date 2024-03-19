package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDCVirtualInterface_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_dc_virtual_interfaces.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcDirectConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDCVirtualInterfaces_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "virtual_interfaces.0.id"),
					resource.TestCheckResourceAttrSet(rName, "virtual_interfaces.0.name"),
					resource.TestCheckResourceAttrSet(rName, "virtual_interfaces.0.status"),
					resource.TestCheckResourceAttrSet(rName, "virtual_interfaces.0.direct_connect_id"),
					resource.TestCheckResourceAttrSet(rName, "virtual_interfaces.0.vgw_id"),
					resource.TestCheckResourceAttrSet(rName, "virtual_interfaces.0.enterprise_project_id"),

					resource.TestCheckOutput("virtual_interface_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("direct_connect_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vgw_id_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDCVirtualInterfaces_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dc_virtual_interfaces" "test" {
  depends_on = [huaweicloud_dc_virtual_interface.test]
}

data "huaweicloud_dc_virtual_interfaces" "virtual_interface_id_filter" {
  virtual_interface_id = huaweicloud_dc_virtual_interface.test.id
}

locals {
  virtual_interface_id = huaweicloud_dc_virtual_interface.test.id
}

output "virtual_interface_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_interfaces.virtual_interface_id_filter.virtual_interfaces) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_interfaces.virtual_interface_id_filter.virtual_interfaces[*].id : v == 
  local.virtual_interface_id]
  )  
}

data "huaweicloud_dc_virtual_interfaces" "name_filter" {
  name = huaweicloud_dc_virtual_interface.test.name
}

locals {
  name = huaweicloud_dc_virtual_interface.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_interfaces.name_filter.virtual_interfaces) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_interfaces.name_filter.virtual_interfaces[*].name : v == local.name]
  )
}

data "huaweicloud_dc_virtual_interfaces" "status_filter" {
  status = huaweicloud_dc_virtual_interface.test.status
}

locals {
  status = huaweicloud_dc_virtual_interface.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_interfaces.status_filter.virtual_interfaces) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_interfaces.status_filter.virtual_interfaces[*].status : v == local.status]
  )
}

data "huaweicloud_dc_virtual_interfaces" "direct_connect_id_filter" {
	direct_connect_id = huaweicloud_dc_virtual_interface.test.direct_connect_id
  }
  
  locals {
	direct_connect_id = huaweicloud_dc_virtual_interface.test.direct_connect_id
  }
  
  output "direct_connect_id_filter_is_useful" {
	value = length(data.huaweicloud_dc_virtual_interfaces.direct_connect_id_filter.virtual_interfaces) > 0 && alltrue(
	  [for v in data.huaweicloud_dc_virtual_interfaces.direct_connect_id_filter.virtual_interfaces[*].
	direct_connect_id : v == local.direct_connect_id]
	)
  }

data "huaweicloud_dc_virtual_interfaces" "vgw_id_filter" {
  vgw_id = huaweicloud_dc_virtual_interface.test.vgw_id
}

locals {
  vgw_id = huaweicloud_dc_virtual_interface.test.vgw_id
}

output "vgw_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_interfaces.vgw_id_filter.virtual_interfaces) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_interfaces.vgw_id_filter.virtual_interfaces[*].vgw_id : v == local.vgw_id]
  )
}

data "huaweicloud_dc_virtual_interfaces" "enterprise_project_id_filter" {
  enterprise_project_id = huaweicloud_dc_virtual_interface.test.enterprise_project_id
}

locals {
  enterprise_project_id = huaweicloud_dc_virtual_interface.test.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_interfaces.enterprise_project_id_filter.virtual_interfaces) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_interfaces.enterprise_project_id_filter.virtual_interfaces[*].
  enterprise_project_id : v == local.enterprise_project_id]
  )
}
`, testAccVirtualInterface_basic(name, acctest.RandIntRange(1, 3999)))
}
