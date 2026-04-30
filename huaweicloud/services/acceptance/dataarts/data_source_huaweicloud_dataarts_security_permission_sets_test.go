package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityPermissionSets_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_permission_sets.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byParentId   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_parent_id"
		dcByParentId = acceptance.InitDataSourceCheck(byParentId)

		byManagerId   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_manager_id"
		dcByManagerId = acceptance.InitDataSourceCheck(byManagerId)

		byName   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byTypeFilter   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_type_filter"
		dcByTypeFilter = acceptance.InitDataSourceCheck(byTypeFilter)

		byManagerName   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_manager_name"
		dcByManagerName = acceptance.InitDataSourceCheck(byManagerName)

		byManagerType   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_manager_type"
		dcByManagerType = acceptance.InitDataSourceCheck(byManagerType)

		byDatasourceType   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_datasource_type"
		dcByDatasourceType = acceptance.InitDataSourceCheck(byDatasourceType)

		bySyncStatus   = "data.huaweicloud_dataarts_security_permission_sets.filter_by_sync_status"
		dcBySyncStatus = acceptance.InitDataSourceCheck(bySyncStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsManagerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityPermissionSets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestMatchResourceAttr(all, "permission_sets.#", regexp.MustCompile(`^[1-9]([0-9]*)$`)),
					resource.TestCheckResourceAttrSet(all, "permission_sets.0.id"),
					resource.TestCheckResourceAttrSet(all, "permission_sets.0.name"),
					resource.TestCheckResourceAttrSet(all, "permission_sets.0.parent_id"),
					resource.TestCheckResourceAttrSet(all, "permission_sets.0.type"),
					resource.TestCheckResourceAttrSet(all, "permission_sets.0.sync_status"),
					// Filter by parent_id
					dcByParentId.CheckResourceExists(),
					resource.TestCheckOutput("is_parent_id_filter_useful", "true"),
					// Filter by manager_id
					dcByManagerId.CheckResourceExists(),
					resource.TestCheckOutput("is_manager_id_filter_useful", "true"),
					// Filter by name
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by type_filter
					dcByTypeFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					// Filter by manager_name
					dcByManagerName.CheckResourceExists(),
					resource.TestCheckOutput("is_manager_name_filter_useful", "true"),
					// Filter by manager_type
					dcByManagerType.CheckResourceExists(),
					resource.TestCheckOutput("is_manager_type_filter_useful", "true"),
					// Filter by datasource_type
					dcByDatasourceType.CheckResourceExists(),
					resource.TestCheckOutput("is_datasource_type_filter_useful", "true"),
					// Filter by sync_status
					dcBySyncStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_sync_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityPermissionSets_basic_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_permission_set" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  parent_id    = "0"
  manager_name = "%[2]s_manager"
  manager_id   = "%[3]s"
  manager_type = "USER"
}

resource "huaweicloud_dataarts_security_permission_set" "sub_test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_sub"
  parent_id    = huaweicloud_dataarts_security_permission_set.test.id
  manager_id   = "%[3]s"

  depends_on = [huaweicloud_dataarts_security_permission_set.test]
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_MANAGER_ID)
}

func testAccDataSecurityPermissionSets_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all permission sets and without any filter
data "huaweicloud_dataarts_security_permission_sets" "all" {
  workspace_id = "%[2]s"

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

# Filter by parent_id (root permission sets)
locals {
  parent_id = "0"
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_parent_id" {
  workspace_id = "%[2]s"
  parent_id    = local.parent_id

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  parent_id_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_parent_id.permission_sets[*].parent_id :
      v == local.parent_id
  ]
}

output "is_parent_id_filter_useful" {
  value = length(local.parent_id_filter_result) > 0 && alltrue(local.parent_id_filter_result)
}

# Filter by manager_id
locals {
  manager_id = "%[3]s"
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_manager_id" {
  workspace_id = "%[2]s"
  manager_id   = local.manager_id

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  manager_id_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_manager_id.permission_sets[*].manager_id :
      v == local.manager_id
  ]
}

output "is_manager_id_filter_useful" {
  value = length(local.manager_id_filter_result) > 0 && alltrue(local.manager_id_filter_result)
}

# Filter by name (fuzzy matching)
locals {
  permission_set_name = huaweicloud_dataarts_security_permission_set.test.name
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = local.permission_set_name

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_name.permission_sets[*].name :
      strcontains(v, local.permission_set_name)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by type_filter
locals {
  type_filter = "SUB_PERMISSION_SET"
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_type_filter" {
  workspace_id = "%[2]s"
  type_filter  = local.type_filter

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_type_filter.permission_sets[*].parent_id :
      v != "0"
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by manager_name
locals {
  manager_name = huaweicloud_dataarts_security_permission_set.test.manager_name
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_manager_name" {
  workspace_id = "%[2]s"
  manager_name = local.manager_name

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  manager_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_manager_name.permission_sets[*].manager_name :
      v == local.manager_name
  ]
}

output "is_manager_name_filter_useful" {
  value = length(local.manager_name_filter_result) > 0 && alltrue(local.manager_name_filter_result)
}

# Filter by manager_type
locals {
  manager_type = huaweicloud_dataarts_security_permission_set.test.manager_type
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_manager_type" {
  workspace_id = "%[2]s"
  manager_type = local.manager_type

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  manager_type_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_manager_type.permission_sets[*].manager_type :
      v == local.manager_type
  ]
}

output "is_manager_type_filter_useful" {
  value = length(local.manager_type_filter_result) > 0 && alltrue(local.manager_type_filter_result)
}

# Filter by datasource_type
locals {
  datasource_type = data.huaweicloud_dataarts_security_permission_sets.all.permission_sets[0].datasource_type
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_datasource_type" {
  workspace_id    = "%[2]s"
  datasource_type = local.datasource_type

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  datasource_type_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_datasource_type.permission_sets[*].datasource_type :
      v == local.datasource_type
  ]
}

output "is_datasource_type_filter_useful" {
  value = length(local.datasource_type_filter_result) > 0 && alltrue(local.datasource_type_filter_result)
}

# Filter by sync_status
locals {
  sync_status = data.huaweicloud_dataarts_security_permission_sets.all.permission_sets[0].sync_status
}

data "huaweicloud_dataarts_security_permission_sets" "filter_by_sync_status" {
  workspace_id = "%[2]s"
  sync_status  = local.sync_status

  depends_on = [
    huaweicloud_dataarts_security_permission_set.test,
    huaweicloud_dataarts_security_permission_set.sub_test,
  ]
}

locals {
  sync_status_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_sets.filter_by_sync_status.permission_sets[*].sync_status :
      v == local.sync_status
  ]
}

output "is_sync_status_filter_useful" {
  value = length(local.sync_status_filter_result) > 0 && alltrue(local.sync_status_filter_result)
}
`, testAccDataSecurityPermissionSets_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_MANAGER_ID)
}
