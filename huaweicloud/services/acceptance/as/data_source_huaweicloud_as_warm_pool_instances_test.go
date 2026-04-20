package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWarmPoolInstancesDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_warm_pool_instances.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please configure the AS group ID into the environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWarmPoolInstancesDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "warm_pool_instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "warm_pool_instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "warm_pool_instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "warm_pool_instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "warm_pool_instances.0.status"),
				),
			},
		},
	})
}

func testAccWarmPoolInstancesDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_warm_pool_instances" "test" {
  scaling_group_id = "%s"
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
