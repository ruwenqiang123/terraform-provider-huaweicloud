package workspace

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getResourceAppImageServerFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	return workspace.GetAppImageServerById(client, state.Primary.ID)
}

func getAcceptanceEpsId() string {
	if acceptance.HW_ENTERPRISE_PROJECT_ID_TEST == "" {
		return "0"
	}
	return acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
}

func TestAccResourceAppImageServer_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_image_server.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()

		serverGroup interface{}
		rc          = acceptance.InitResourceCheck(
			resourceName,
			&serverGroup,
			getResourceAppImageServerFunc,
		)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceAppImageSpecCode(t)
			acceptance.TestAccPrecheckWorkspaceUserNames(t)
			acceptance.TestAccPreCheckWorkspaceOUName(t)
			acceptance.TestAccPreCheckWorkspaceADDomainNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceAppImageServer_basic(name, "Created by script"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "authorize_accounts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "authorize_accounts.0.type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "authorize_accounts.0.domain",
						retrieveActiveDomainName(strings.Split(acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES, ","))),
					resource.TestCheckResourceAttr(resourceName, "image_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID),
					resource.TestCheckResourceAttr(resourceName, "image_type", "gold"),
					resource.TestCheckResourceAttr(resourceName, "spec_code",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", getAcceptanceEpsId()),
				),
			},
			{
				Config: testResourceAppImageServer_basic(updateName, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"flavor_id",
					"vpc_id",
					"subnet_id",
					"root_volume",
					"image_source_product_id",
					"is_vdi",
					"availability_zone",
					"ou_name",
					"extra_session_type",
					"extra_session_size",
					"route_policy",
					"scheduler_hints",
					"tags",
					"enterprise_project_id",
					"is_delete_associated_resources",
				},
			},
		},
	})
}

func testResourceAppImageServer_basic(name, description string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_app_image_server" "test" {
  name                    = "%[1]s"
  flavor_id               = "%[2]s"
  vpc_id                  = "%[3]s"
  subnet_id               = "%[4]s"
  image_id                = "%[5]s"
  image_type              = "gold"
  image_source_product_id = "%[6]s"
  spec_code               = "%[7]s"

  # Currently only one user can be set.
  authorize_accounts {
    account = split(",", "%[8]s")[0]
    type    = "USER"
    domain  = element(split(",", "%[9]s"), 0)
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  is_vdi             = false
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  description        = "%[10]s"
  ou_name            = "%[11]s"
  extra_session_type = "CPU"
  extra_session_size = 2

  route_policy {
    max_session   = 3
    cpu_threshold = 80
    mem_threshold = 80
  }

  tags = {
    foo = "bar"
  }

  enterprise_project_id          = "%[12]s"
  is_delete_associated_resources = true
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE,
		acceptance.HW_WORKSPACE_USER_NAMES,
		acceptance.HW_WORKSPACE_AD_DOMAIN_NAMES,
		description,
		acceptance.HW_WORKSPACE_OU_NAME,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
	)
}
