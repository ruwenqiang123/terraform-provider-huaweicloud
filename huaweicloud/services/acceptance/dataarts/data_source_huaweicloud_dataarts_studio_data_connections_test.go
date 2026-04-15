package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataStudioDataConnections_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_studio_data_connections.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byConnectionId   = "data.huaweicloud_dataarts_studio_data_connections.filter_by_connection_id"
		dcByConnectionId = acceptance.InitDataSourceCheck(byConnectionId)

		byName   = "data.huaweicloud_dataarts_studio_data_connections.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_dataarts_studio_data_connections.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataStudioDataConnections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "connections.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByConnectionId.CheckResourceExists(),
					resource.TestCheckOutput("is_connection_id_filter_useful", "true"),
					resource.TestCheckResourceAttr(byConnectionId, "connections.#", "1"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.id"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.name"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.type"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.qualified_name"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.created_by"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.created_at"),
					resource.TestCheckResourceAttrSet(byConnectionId, "connections.0.create_timestamp"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataStudioDataConnections_basic() string {
	name := acceptance.RandomAccResourceName()

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
    "cdm_property_enable": "false"
  })

  lifecycle {
    ignore_changes = [
      config,
    ]
  }
}

# Query all data connections without any filter.
data "huaweicloud_dataarts_studio_data_connections" "all" {
  depends_on = [huaweicloud_dataarts_studio_data_connection.test]

  workspace_id = "%[2]s"
}

# Filter by connection ID.
locals {
  data_connection_id = var.data_connection_id != "" ? var.data_connection_id : huaweicloud_dataarts_studio_data_connection.test[0].id
}

data "huaweicloud_dataarts_studio_data_connections" "filter_by_connection_id" {
  workspace_id  = "%[2]s"
  connection_id = local.data_connection_id
}

locals {
  connection_id_filter_result = [
    for v in data.huaweicloud_dataarts_studio_data_connections.filter_by_connection_id.connections[*].id : v == local.data_connection_id
  ]
}

output "is_connection_id_filter_useful" {
  value = length(local.connection_id_filter_result) > 0 && alltrue(local.connection_id_filter_result)
}

# Filter by name.
locals {
  data_connection_name = var.data_connection_id != "" ? try(data.huaweicloud_dataarts_studio_data_connections.all.connections[0].name,
    "NOT_FOUND") : "%[3]s"
}

data "huaweicloud_dataarts_studio_data_connections" "filter_by_name" {
  depends_on   = [huaweicloud_dataarts_studio_data_connection.test]

  workspace_id = "%[2]s"
  name         = local.data_connection_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_studio_data_connections.filter_by_name.connections[*].name : strcontains(v, local.data_connection_name)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}


# Filter by type.
locals {
  data_connection_type = var.data_connection_id != "" ? try(data.huaweicloud_dataarts_studio_data_connections.all.connections[0].type,
    "NOT_FOUND") : huaweicloud_dataarts_studio_data_connection.test[0].type
}

data "huaweicloud_dataarts_studio_data_connections" "filter_by_type" {
  depends_on   = [huaweicloud_dataarts_studio_data_connection.test]

  workspace_id = "%[2]s"
  type         = local.data_connection_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dataarts_studio_data_connections.filter_by_type.connections[*].type : v == local.data_connection_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`, acceptance.HW_DATAARTS_CONNECTION_ID, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
