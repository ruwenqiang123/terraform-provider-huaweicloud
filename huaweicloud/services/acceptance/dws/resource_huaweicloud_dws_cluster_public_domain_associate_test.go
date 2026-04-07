package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getClusterAssociatedPublicDomainFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	return dws.GetClusterPublicEndpointById(client, state.Primary.ID)
}

// Before running this test, you need to bind public EIP to at least one DWS cluster node.
func TestAccClusterPublicDomainAssociate_basic(t *testing.T) {
	var (
		obj        interface{}
		rName      = "huaweicloud_dws_cluster_public_domain_associate.test"
		name       = acceptance.RandomAccResourceNameWithDash()
		updateName = acceptance.RandomAccResourceNameWithDash()

		rc = acceptance.InitResourceCheck(rName, &obj, getClusterAssociatedPublicDomainFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testClusterPublicDomainAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "domain_name", name),
					resource.TestMatchResourceAttr(rName, "ttl", regexp.MustCompile("^[1-9][0-9]*$")),
				),
			},
			{
				Config: testClusterPublicDomainAssociate_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", updateName),
					resource.TestCheckResourceAttr(rName, "ttl", "1000"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: clusterPublicDomainAssociateImportState(rName),
			},
		},
	})
}

func clusterPublicDomainAssociateImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.Attributes["domain_name"]), nil
	}
}

func testClusterPublicDomainAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_public_domain_associate" "test" {
  cluster_id  = "%[1]s"
  domain_name = "%[2]s"
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}

func testClusterPublicDomainAssociate_basic_step2(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_public_domain_associate" "test" {
  cluster_id  = "%[1]s"
  domain_name = "%[2]s"
  ttl         = 1000
}
`, acceptance.HW_DWS_CLUSTER_ID, updateName)
}
