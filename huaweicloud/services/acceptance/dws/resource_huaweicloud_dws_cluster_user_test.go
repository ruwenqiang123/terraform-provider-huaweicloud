package dws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getClusterUserResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		clusterId = state.Primary.Attributes["cluster_id"]
		name      = state.Primary.Attributes["name"]
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	return dws.GetClusterUser(client, clusterId, name)
}

func TestAccResourceClusterUser_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_dws_cluster_user.test"
		obj          interface{}
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getClusterUserResourceFunc)

		rName = strings.ToLower(acceptance.RandomAccResourceName())
	)

	const userType = "user"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsGrantTargets(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceClusterUser_user_step1(rName, userType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", userType),
					resource.TestCheckResourceAttr(resourceName, "login", "true"),
					resource.TestCheckResourceAttr(resourceName, "create_db", "false"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
				),
			},
			{
				Config: testAccResourceClusterUser_user_step2(rName, userType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "login", "false"),
					resource.TestCheckResourceAttr(resourceName, "create_db", "true"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceClusterUserImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"password",
					"type",
					"cascade",
				},
			},
		},
	})
}

func testAccResourceClusterUser_user_step1(name, userType string) string {
	return testAccResourceClusterUser_basic(name, userType)
}

func testAccResourceClusterUser_user_step2(name, userType string) string {
	return testAccResourceClusterUser_update(name, userType)
}

func TestAccResourceClusterUser_role(t *testing.T) {
	var (
		resourceName = "huaweicloud_dws_cluster_user.test"
		obj          interface{}
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getClusterUserResourceFunc)

		rName = strings.ToLower(acceptance.RandomAccResourceName())
	)

	const userType = "role"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsGrantTargets(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceClusterUser_role_step1(rName, userType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", userType),
					resource.TestCheckResourceAttr(resourceName, "login", "true"),
					resource.TestCheckResourceAttr(resourceName, "create_db", "false"),
				),
			},
			{
				Config: testAccResourceClusterUser_role_step2(rName, userType),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "login", "false"),
					resource.TestCheckResourceAttr(resourceName, "create_db", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceClusterUserImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"password",
					"type",
					"cascade",
				},
			},
		},
	})
}

func testAccResourceClusterUser_role_step1(name, userType string) string {
	return testAccResourceClusterUser_basic(name, userType)
}

func testAccResourceClusterUser_role_step2(name, userType string) string {
	return testAccResourceClusterUser_update(name, userType)
}

func testAccResourceClusterUser_basic(name, userType string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_user" "test" {
  cluster_id  = "%[1]s"
  name        = "%[2]s"
  type        = "%[3]s"
  password    = "HuaweiTest@123456789"
  cascade     = false
  login       = true
  create_role = false
  create_db   = false
  inherit     = true
  conn_limit  = -1
}
`, acceptance.HW_DWS_CLUSTER_ID, name, userType)
}

func testAccResourceClusterUser_update(name, userType string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_cluster_user" "test" {
  cluster_id  = "%[1]s"
  name        = "%[2]s"
  type        = "%[3]s"
  password    = "HuaweiTest@123456789"
  cascade     = true
  login       = false
  create_role = false
  create_db   = true
  inherit     = true
  conn_limit  = -1
}
`, acceptance.HW_DWS_CLUSTER_ID, name, userType)
}

func testAccResourceClusterUserImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterId, name string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_dws_cluster_user" {
				clusterId = rs.Primary.Attributes["cluster_id"]
				name = rs.Primary.Attributes["name"]
			}
		}
		if clusterId == "" || name == "" {
			return "", fmt.Errorf("resource not found: %s/%s", clusterId, name)
		}

		return fmt.Sprintf("%s/%s/%s", acceptance.HW_REGION_NAME, clusterId, name), nil
	}
}

func TestAccResourceClusterUser_advanced(t *testing.T) {
	var (
		rName = strings.ToLower(acceptance.RandomAccResourceName())

		obj                interface{}
		grantsResourceName = "huaweicloud_dws_cluster_user.test"
		rcGrants           = acceptance.InitResourceCheck(
			grantsResourceName,
			&obj,
			getClusterUserResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcGrants.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceClusterUser_withGrants(rName),
				Check: resource.ComposeTestCheckFunc(
					rcGrants.CheckResourceExists(),
					resource.TestCheckResourceAttr(grantsResourceName, "name", fmt.Sprintf("%s_grants", rName)),
					resource.TestCheckResourceAttr(grantsResourceName, "type", "user"),
					resource.TestCheckResourceAttr(grantsResourceName, "grant_list.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(grantsResourceName, "grant_list.*", map[string]string{
						"type":         "DATABASE",
						"database":     acceptance.HW_DWS_GRANT_DATABASE_NAME,
						"privileges.#": "2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(grantsResourceName, "grant_list.*", map[string]string{
						"type":         "TABLE",
						"database":     acceptance.HW_DWS_GRANT_DATABASE_NAME,
						"schema_name":  acceptance.HW_DWS_GRANT_SCHEMA_NAME,
						"object_name":  acceptance.HW_DWS_GRANT_OBJECT_NAME,
						"privileges.#": "2",
					}),
				),
			},
		},
	})
}

func testAccResourceClusterUser_withGrants(name string) string {
	return fmt.Sprintf(`
variable "grants" {
  type = list(object({
    type        = string
    database    = string
    schema_name = optional(string)
    object_name = optional(string)

    privileges = list(object({
      permission = string
      grant_with = bool
    }))
  }))

  default = [
    {
      type     = "DATABASE"
      database = "%[1]s"

      privileges = [
        {
          permission = "CONNECT"
          grant_with = false
        },
        {
          permission = "CREATE"
          grant_with = false
        }
      ]
    },
    {
      type        = "TABLE"
      database    = "%[1]s"
      schema_name = "%[2]s"
      object_name = "%[3]s"

      privileges = [
        {
          permission = "SELECT"
          grant_with = false
        },
        {
          permission = "INSERT"
          grant_with = false
        }
      ]
    }
  ]
}

resource "huaweicloud_dws_cluster_user" "test" {
  cluster_id = "%[4]s"
  name       = "%[5]s_grants"
  type       = "user"
  password   = "HuaweiTest@123456789"

  dynamic "grant_list" {
    for_each = var.grants

    content {
      type        = grant_list.value.type
      database    = grant_list.value.database
      schema_name = lookup(grant_list.value, "schema_name", null)
      object_name = lookup(grant_list.value, "object_name", null)

      dynamic "privileges" {
        for_each = grant_list.value.privileges

        content {
          permission = privileges.value.permission
          grant_with = privileges.value.grant_with
        }
      }
    }
  }
}`, acceptance.HW_DWS_GRANT_DATABASE_NAME,
		acceptance.HW_DWS_GRANT_SCHEMA_NAME,
		acceptance.HW_DWS_GRANT_OBJECT_NAME,
		acceptance.HW_DWS_CLUSTER_ID,
		name)
}
