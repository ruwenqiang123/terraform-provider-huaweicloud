package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDataServiceApprovers_basic(t *testing.T) {
	var (
		name = "data.huaweicloud_dataarts_dataservice_approvers.test"
		dc   = acceptance.InitDataSourceCheck(name)

		byName   = "data.huaweicloud_dataarts_dataservice_approvers.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDataServiceApprovers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(name, "approvers.#", regexp.MustCompile(`^\d+$`)),
					resource.TestCheckResourceAttrSet(name, "approvers.0.id"),
					resource.TestCheckResourceAttrSet(name, "approvers.0.name"),
					resource.TestCheckResourceAttrSet(name, "approvers.0.user_id"),
					resource.TestCheckResourceAttrSet(name, "approvers.0.user_name"),
					resource.TestCheckResourceAttrSet(name, "approvers.0.create_by"),
					resource.TestMatchResourceAttr(name, "approvers.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(name, "approvers.0.approver_type", regexp.MustCompile(`^\d+$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDataServiceApprovers_basic() string {
	return fmt.Sprintf(`
# Query all approvers and without any filter
data "huaweicloud_dataarts_dataservice_approvers" "test" {
  workspace_id = "%[1]s"
}

# Filter by name
locals {
  approver_name = try(data.huaweicloud_dataarts_dataservice_approvers.test.approvers[0].name, "NOT_FOUND")
}

data "huaweicloud_dataarts_dataservice_approvers" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.approver_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_dataservice_approvers.filter_by_name.approvers[*].name 
    : v == local.approver_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
