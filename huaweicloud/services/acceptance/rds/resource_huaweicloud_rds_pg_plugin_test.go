package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPgPluginResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getPgPluginClient, err := cfg.NewServiceClient("rds", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	listPgPluginHttpUrl := "v3/{project_id}/instances/{instance_id}/extensions?database_name={database_name}"
	listPgPluginPath := getPgPluginClient.Endpoint + listPgPluginHttpUrl
	listPgPluginPath = strings.ReplaceAll(listPgPluginPath, "{project_id}", getPgPluginClient.ProjectID)
	listPgPluginPath = strings.ReplaceAll(listPgPluginPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	listPgPluginPath = strings.ReplaceAll(listPgPluginPath, "{database_name}", state.Primary.Attributes["database_name"])

	resp, err := getPgPluginClient.Request("GET", listPgPluginPath, &golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL : %s", err)
	}

	body, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL : %s", err)
	}

	name := state.Primary.Attributes["name"]
	plugin := utils.PathSearch(fmt.Sprintf("extensions[?name=='%s']|[?created]|[0]", name), body, nil)

	if plugin == nil {
		return nil, fmt.Errorf("no RDS PostgreSQL plugin matching %s was found", name)
	}

	return plugin, nil
}

func TestAccPgPlugin_basic(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resourceName := "huaweicloud_rds_pg_plugin.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPgPluginResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgPlugin_basic(randName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "pgl_ddl_deploy"),
					resource.TestCheckResourceAttr(resourceName, "database_name", "postgres"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
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

func testPgPlugin_basic(randName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_pg_plugin" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "pgl_ddl_deploy"
  database_name = "postgres"
}
`, testAccRdsInstance_basic(randName, pwd))
}
