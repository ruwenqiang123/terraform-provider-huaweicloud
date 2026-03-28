package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before testing, please ensure that you have created an OBS bucket and set the root directory of the bucket to the log
// storage path of your DataArts workspace.
func TestAccFactoryJobExport_basic(t *testing.T) {
	var (
		name                  = acceptance.RandomAccResourceName()
		withExportDependTrue  = "huaweicloud_dataarts_factory_job_export.with_export_depend_true"
		withExportDependFalse = "huaweicloud_dataarts_factory_job_export.with_export_depend_false"
		withoutExportDepend   = "huaweicloud_dataarts_factory_job_export.without_export_depend"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsCdmName(t)
			acceptance.TestAccPreCheckOBSObjectStoragePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFactoryJobExport_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(withExportDependTrue, "folder_path"),
					resource.TestCheckResourceAttrSet(withExportDependFalse, "folder_path"),
					resource.TestCheckResourceAttrSet(withoutExportDepend, "folder_path"),
				),
			},
		},
	})
}

func testAccFactoryJobExport_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_job" "test" {
  name         = "%[1]s"
  workspace_id = "%[2]s"
  process_type = "BATCH"

  nodes {
    name = "Rest_client_%[1]s"
    type = "RESTAPI"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "url"
      value = "https://www.huaweicloud.com/"
    }

    properties {
      name  = "method"
      value = "GET"
    }

    properties {
      name  = "retry"
      value = "false"
    }

    properties {
      name  = "requestMode"
      value = "sync"
    }

    properties {
      name  = "securityAuthentication"
      value = "NONE"
    }

    properties {
      name  = "agentName"
      value = "%[3]s" # The agentName obtained from the DataArts Migration side.
    }
  }

  schedule {
    type = "EXECUTE_ONCE"
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CDM_NAME)
}

func testAccFactoryJobExport_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_factory_job_export" "with_export_depend_true" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test
  ]

  workspace_id  = "%[2]s"
  job_name      = huaweicloud_dataarts_factory_job.test.name
  export_depend = true
  export_status = "DEVELOP"
  obs_path      = "%[3]s" # The OBS path for storing the exported job package, such as obs://dataarts-test/jobs-storage/
}

resource "huaweicloud_dataarts_factory_job_export" "with_export_depend_false" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test
  ]

  workspace_id  = "%[2]s"
  job_name      = huaweicloud_dataarts_factory_job.test.name
  export_depend = false
  export_status = "DEVELOP"
  obs_path      = "%[3]s" # The OBS path for storing the exported job package, such as obs://dataarts-test/jobs-storage/
}

resource "huaweicloud_dataarts_factory_job_export" "without_export_depend" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test
  ]

  workspace_id  = "%[2]s"
  job_name      = huaweicloud_dataarts_factory_job.test.name
  export_status = "DEVELOP"
  obs_path      = "%[3]s" # The OBS path for storing the exported job package, such as obs://dataarts-test/jobs-storage/
}
`, testAccFactoryJobExport_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_OBS_OBJECT_STORAGE_PATH)
}
