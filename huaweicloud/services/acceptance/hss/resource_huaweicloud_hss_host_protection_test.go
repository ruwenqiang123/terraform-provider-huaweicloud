package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/hss"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getHostProtectionFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcHssV5Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HSS v5 client: %s", err)
	}

	var (
		epsId = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		id    = state.Primary.ID
	)

	// If the enterprise project ID is not set during query, query all enterprise projects.
	if epsId == "" {
		epsId = hss.QueryAllEpsValue
	}
	listHostOpts := hssv5model.ListHostStatusRequest{
		Region:              &acceptance.HW_REGION_NAME,
		EnterpriseProjectId: utils.String(epsId),
		HostId:              utils.String(id),
	}

	resp, err := client.ListHostStatus(&listHostOpts)
	if err != nil {
		return nil, fmt.Errorf("error querying HSS hosts: %s", err)
	}

	if resp == nil || resp.DataList == nil {
		return nil, fmt.Errorf("the host (%s) for HSS host protection does not exist", id)
	}

	hostList := *resp.DataList
	if len(hostList) == 0 || utils.StringValue(hostList[0].ProtectStatus) == string(hss.ProtectStatusClosed) {
		return nil, golangsdk.ErrDefault404{}
	}

	return hostList[0], nil
}

func TestAccHostProtection_basic(t *testing.T) {
	var (
		host  *hssv5model.Host
		rName = "huaweicloud_hss_host_protection.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&host,
		getHostProtectionFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
			acceptance.TestAccPreCheckHSSHostProtectionQuotaId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccHostProtection_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttr(rName, "version", "hss.version.basic"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(rName, "quota_id", acceptance.HW_HSS_HOST_PROTECTION_QUOTA_ID),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "host_status"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "agent_id"),
					resource.TestCheckResourceAttrSet(rName, "agent_status"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "detect_result"),
					resource.TestCheckResourceAttrSet(rName, "asset_value"),
					resource.TestCheckResourceAttrSet(rName, "open_time"),
				),
			},
			{
				Config: testAccHostProtection_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "host_id", acceptance.HW_HSS_HOST_PROTECTION_HOST_ID),
					resource.TestCheckResourceAttr(rName, "version", "hss.version.enterprise"),
					resource.TestCheckResourceAttr(rName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "host_name"),
					resource.TestCheckResourceAttrSet(rName, "host_status"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "agent_id"),
					resource.TestCheckResourceAttrSet(rName, "agent_status"),
					resource.TestCheckResourceAttrSet(rName, "os_type"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "detect_result"),
					resource.TestCheckResourceAttrSet(rName, "asset_value"),
					resource.TestCheckResourceAttrSet(rName, "open_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"quota_id",
				},
			},
		},
	})
}

func testAccHostProtection_basic() string {
	return fmt.Sprintf(`

resource "huaweicloud_hss_host_protection" "test" {
  host_id               = "%[1]s"
  version               = "hss.version.basic"
  charging_mode         = "prePaid"
  quota_id              = "%[2]s"
  enterprise_project_id = "%[3]s"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID, acceptance.HW_HSS_HOST_PROTECTION_QUOTA_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccHostProtection_update() string {
	return fmt.Sprintf(`

resource "huaweicloud_hss_host_protection" "test" {
  host_id               = "%[1]s"
  version               = "hss.version.enterprise"
  charging_mode         = "postPaid"
  enterprise_project_id = "%[2]s"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
