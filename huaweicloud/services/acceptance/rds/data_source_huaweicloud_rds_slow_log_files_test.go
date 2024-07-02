package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSlowLogFiles_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_slow_log_files.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsSlowLogFiles_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "files.#"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.file_name"),
					resource.TestCheckResourceAttrSet(dataSource, "files.0.file_size"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsSlowLogFiles_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_slow_log_files" "test" {
  instance_id = "%s"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
