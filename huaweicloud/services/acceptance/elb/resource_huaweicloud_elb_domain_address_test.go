package elb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
)

func getDomainAddressResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("elb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	resourceId := strings.Split(state.Primary.ID, "/")

	return elb.GetDomainAddressInfomation(client, resourceId[0], resourceId[1])
}

func TestAccDomainAddress_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_elb_domain_address.test"
		obj          interface{}
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDomainAddressResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Beforce running test, prepare a loadbalancer and enable domain resolution (private and public).
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainAddressConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "loadbalancer_id", acceptance.HW_ELB_LOADBALANCER_ID),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "enable"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDomainAddressConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_elb_loadbalancers" "test" {
  loadbalancer_id = "%[1]s"
}

resource "huaweicloud_elb_domain_address" "test" {
  loadbalancer_id = "%[1]s"
  ip_address      = data.huaweicloud_elb_loadbalancers.test.loadbalancers[0].ipv4_address
}
`, acceptance.HW_ELB_LOADBALANCER_ID)
}
