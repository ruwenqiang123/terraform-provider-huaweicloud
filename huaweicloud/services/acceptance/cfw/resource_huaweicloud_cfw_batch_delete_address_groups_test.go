package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchDeleteAddressGroups_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwAddressGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBatchDeleteAddressGroups_basic(),
			},
		},
	})
}

func testBatchDeleteAddressGroups_base() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

func testBatchDeleteAddressGroups_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_batch_delete_address_groups" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  set_ids   = ["%[2]s"]
}
`, testBatchDeleteAddressGroups_base(), acceptance.HW_CFW_ADDRESS_GROUP_ID)
}
