package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAgencyPermission_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAgencyPermission_basic(),
			},
		},
	})
}

func testAgencyPermission_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_css_agency_permission" "test" {
  domain_id   = "%[1]s"
  domain_name = "%[2]s"
  type        = "vpc"
}
`, acceptance.HW_DOMAIN_ID, acceptance.HW_DOMAIN_NAME)
}
