package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBackupStop_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBackupStop_basic(),
			},
		},
	})
}

func testBackupStop_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_backups" "test" {
  instance_id = "%s"
}

resource "huaweicloud_dds_backup_stop" "test" {
  backup_id = data.huaweicloud_dds_backups.test.backups[0].id
  action    = "stop"
}
`, acceptance.HW_DDS_INSTANCE_ID)
}
