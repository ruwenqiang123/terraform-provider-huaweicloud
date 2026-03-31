package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFactoryJobImport_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

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
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFactoryJobImport_basic(name),
			},
		},
	})
}

func testAccFactoryJobImport_base(name string) string {
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

resource "huaweicloud_dataarts_factory_job_export" "test" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test
  ]

  workspace_id  = "%[2]s"
  job_name      = huaweicloud_dataarts_factory_job.test.name
  export_status = "DEVELOP"
  obs_path      = "%[4]s" # The OBS path for storing the exported job package, such as obs://dataarts-test/jobs-storage/
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CDM_NAME, acceptance.HW_OBS_OBJECT_STORAGE_PATH)
}

func testAccFactoryJobImport_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "time_sleep" "test" {
  depends_on = [
    huaweicloud_dataarts_factory_job_export.test
  ]

  create_duration = "10s"
}

resource "huaweicloud_dataarts_factory_job_import" "test" {
  depends_on = [
    time_sleep.test
  ]

  workspace_id     = "%[2]s"
  path             = format("%[3]s%%s/%[4]s.zip", huaweicloud_dataarts_factory_job_export.test.folder_path)
  same_name_policy = "OVERWRITE"
  target_status    = "SAVED"
}
`, testAccFactoryJobImport_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_OBS_OBJECT_STORAGE_PATH,
		name)
}
