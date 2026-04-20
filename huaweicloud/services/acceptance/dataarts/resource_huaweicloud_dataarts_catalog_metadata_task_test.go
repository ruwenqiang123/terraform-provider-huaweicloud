package dataarts

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getCatalogMetadataTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts client: %s", err)
	}
	return dataarts.GetCatalogMetadataTaskById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccCatalogMetadataTask_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dataarts_catalog_metadata_task.test"
		rcCreateTask = acceptance.InitResourceCheck(resourceName, &obj, getCatalogMetadataTaskResourceFunc)

		name        = acceptance.RandomAccResourceName()
		updateName  = acceptance.RandomAccResourceName()
		currentTime = time.Now().Local().Format(time.RFC3339)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcCreateTask.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCatalogMetadataTask_basic_step1(name, currentTime),
				Check: resource.ComposeTestCheckFunc(
					rcCreateTask.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dir_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.cron_expression", "00 */15 9-23 * * ?"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.max_time_out", "60"),
					resource.TestCheckResourceAttrSet(resourceName, "schedule_config.0.end_time"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.interval", "15 minutes"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.schedule_type", "CRON"),
					resource.TestCheckResourceAttrSet(resourceName, "schedule_config.0.start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "schedule_config.0.job_id"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.enabled", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_source_type", "DLI"),
					resource.TestCheckResourceAttrSet(resourceName, "task_config"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "terminal_before_modify", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", ""),
					resource.TestCheckResourceAttrSet(resourceName, "duty_person"),
				),
			},
			{
				Config: testAccCatalogMetadataTask_basic_step2(updateName, currentTime),
				Check: resource.ComposeTestCheckFunc(
					rcCreateTask.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "dir_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.cron_expression", "00 */16 8-23 * * ?"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.max_time_out", "120"),
					resource.TestCheckResourceAttrSet(resourceName, "schedule_config.0.end_time"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.interval", "30 minutes"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.schedule_type", "CRON"),
					resource.TestCheckResourceAttrSet(resourceName, "schedule_config.0.start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "schedule_config.0.job_id"),
					resource.TestCheckResourceAttr(resourceName, "schedule_config.0.enabled", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_source_type", "DLI"),
					resource.TestCheckResourceAttrSet(resourceName, "task_config"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "terminal_before_modify", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "update_time"),
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttr(resourceName, "path", ""),
					resource.TestCheckResourceAttrSet(resourceName, "duty_person"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCatalogMetadataTaskImportStateFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"terminal_before_modify",
				},
			},
		},
	})
}

func testAccCatalogMetadataTask_basic_step1(name, currentTime string) string {
	return fmt.Sprintf(`
variable "data_connection_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_dataarts_studio_data_connection" "test" {
  count = var.data_connection_id == "" ? 1 : 0

  workspace_id = "%[2]s"
  type         = "DLI"
  name         = "%[3]s"
  env_type     = 0
  config       = jsonencode({
    "cdmPropertyEnable"             = false
    "metadata.collectionScope"      = ""
    "metadata.enableAutoCollection" = false
    "metadata.enableRealtimeSync"   = false
  })

  lifecycle {
    ignore_changes = [
      config,
    ]
  }
}

data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = "%[2]s"
  connection_id = var.data_connection_id != "" ? var.data_connection_id : huaweicloud_dataarts_studio_data_connection.test[0].id
}

resource "huaweicloud_dataarts_catalog_metadata_task" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s"
  dir_id       = "0"

  schedule_config {
    cron_expression = "00 */15 9-23 * * ?"
    max_time_out    = 60
    end_time        = format("%%s +08", split("+", timeadd("%[4]s", "24h"))[0])
    interval        = "15 minutes"
    schedule_type   = "CRON"
    start_time      = format("%%s +08", split("+", timeadd("%[4]s", "1h"))[0])
    enabled         = 1
  }

  data_source_type = "DLI"
  task_config      = jsonencode({
    data_connection_name        = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, "")
    data_connection_id          = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].id, "")
    data_connection_create_time = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].create_timestamp, "")
    databaseName                = [
      "tf_test_randx",
      "tpch"
    ]
    tableName                   = [
      "tf_test_randx.tf_test_randx",
      "tpch.part",
      "tpch.region",
      "tpch.customer"
    ]
    alive_object_policy         = "3"
    deleted_obkect_policy       = "1"
    deleted_object_policy       = "10"
    enableDataProfile           = true
    enableDataClassification    = false
    sampling                    = "10"
    queue                       = "default"
    unique                      = true
  })
  description      = "Created by terraform script"
}
`, acceptance.HW_DATAARTS_CONNECTION_ID, acceptance.HW_DATAARTS_WORKSPACE_ID, name, currentTime)
}

func testAccCatalogMetadataTask_basic_step2(name, currentTime string) string {
	return fmt.Sprintf(`
variable "data_connection_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_dataarts_studio_data_connection" "test" {
  count = var.data_connection_id == "" ? 1 : 0

  workspace_id = "%[2]s"
  type         = "DLI"
  name         = "%[3]s"
  env_type     = 0
  config       = jsonencode({
    "cdmPropertyEnable"             = false
    "metadata.collectionScope"      = ""
    "metadata.enableAutoCollection" = false
    "metadata.enableRealtimeSync"   = false
  })

  lifecycle {
    ignore_changes = [
      config,
    ]
  }
}

data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = "%[2]s"
  connection_id = var.data_connection_id != "" ? var.data_connection_id : huaweicloud_dataarts_studio_data_connection.test[0].id
}

resource "huaweicloud_dataarts_catalog_metadata_task" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s"
  dir_id       = "0"

  schedule_config {
    cron_expression = "00 */16 8-23 * * ?"
    max_time_out    = 120
    end_time        = format("%%s +08", split("+", timeadd("%[4]s", "48h"))[0])
    interval        = "30 minutes"
    schedule_type   = "CRON"
    start_time      = format("%%s +08", split("+", timeadd("%[4]s", "2h"))[0])
    enabled         = 1
  }

  data_source_type       = "DLI"
  task_config            = jsonencode({
    data_connection_name      = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, "")
    data_connection_id        = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].id, "")
    databaseName              = [
      "tpch"
    ]
    tableName                 = [
      "tpch.part",
      "tpch.region",
    ]
    alive_object_policy       = "4"
    deleted_obkect_policy     = "2"
    deleted_object_policy     = "11"
    enableDataProfile         = true
    enableDataClassification  = false
    sampling                  = "11"
    queue                     = "default"
    unique                    = true
  })
  terminal_before_modify = true
}
`, acceptance.HW_DATAARTS_CONNECTION_ID, acceptance.HW_DATAARTS_WORKSPACE_ID, name, currentTime)
}

func testAccCatalogMetadataTaskImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		taskId := rs.Primary.ID
		if workspaceId == "" || taskId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, taskId)
		}
		return fmt.Sprintf("%s/%s", workspaceId, taskId), nil
	}
}
