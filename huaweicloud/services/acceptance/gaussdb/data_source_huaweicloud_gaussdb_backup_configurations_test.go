package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBackupConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_backup_configurations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBackupConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rate_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "file_split_size"),
					resource.TestCheckResourceAttrSet(dataSource, "prefetch_block"),
					resource.TestCheckResourceAttrSet(dataSource, "enable_standby_backup"),
					resource.TestCheckResourceAttrSet(dataSource, "close_compression"),
					resource.TestCheckResourceAttrSet(dataSource, "default_backup_method"),
					resource.TestCheckResourceAttrSet(dataSource, "default_backup_media_type"),
				),
			},
		},
	})
}

func testDataSourceBackupConfigurations_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_backup_configurations" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
