package dws

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dws"
)

func getClusterExceptionRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS Client: %s", err)
	}

	return dws.GetClusterExceptionRuleConfigurations(client, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["name"], nil, false)
}

func TestAccClusterExceptionRule_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dws_cluster_exception_rule.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getClusterExceptionRuleFunc)

		name = acceptance.RandomAccResourceName()
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
				Config: testAccClusterExceptionRule_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "configurations.#", "6"),
					resource.TestCheckResourceAttr(rName, "configurations.0.key", "action"),
					resource.TestCheckResourceAttr(rName, "configurations.0.value", "abort"),
					resource.TestCheckResourceAttr(rName, "configurations.1.key", "blocktime"),
					resource.TestCheckResourceAttr(rName, "configurations.1.value", "300"),
					resource.TestCheckResourceAttr(rName, "configurations.2.key", "elapsedtime"),
					resource.TestCheckResourceAttr(rName, "configurations.2.value", "400"),
					resource.TestCheckResourceAttr(rName, "configurations.3.key", "allcputime"),
					resource.TestCheckResourceAttr(rName, "configurations.3.value", "500"),
					resource.TestCheckResourceAttr(rName, "configurations.4.key", "cpuskewpercent"),
					resource.TestCheckResourceAttr(rName, "configurations.4.value", "60"),
					resource.TestCheckResourceAttr(rName, "configurations.5.key", "cpuavgpercent"),
					resource.TestCheckResourceAttr(rName, "configurations.5.value", "70"),
				),
			},
			{
				Config: testAccClusterExceptionRule_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "configurations.#", "5"),
					resource.TestCheckResourceAttr(rName, "configurations.0.key", "action"),
					resource.TestCheckResourceAttr(rName, "configurations.0.value", "penalty"),
					resource.TestCheckResourceAttr(rName, "configurations.1.key", "spillsize"),
					resource.TestCheckResourceAttr(rName, "configurations.1.value", "300"),
					resource.TestCheckResourceAttr(rName, "configurations.2.key", "bandwidth"),
					resource.TestCheckResourceAttr(rName, "configurations.2.value", "400"),
					resource.TestCheckResourceAttr(rName, "configurations.3.key", "memsize"),
					resource.TestCheckResourceAttr(rName, "configurations.3.value", "500"),
					resource.TestCheckResourceAttr(rName, "configurations.4.key", "broadcastsize"),
					resource.TestCheckResourceAttr(rName, "configurations.4.value", "600"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"configurations",
					"configurations_origin",
				},
				ImportStateCheck: checkClusterExceptionRuleImportState,
			},
		},
	})
}

// checkClusterExceptionRuleImportState Verifies the exception rule configuration after import
// 1. The configuration in step two must exist and have the correct value
// 2. The old configuration in step one must be -1
func checkClusterExceptionRuleImportState(states []*terraform.InstanceState) error {
	if len(states) == 0 {
		return errors.New("no state after import")
	}
	is := states[0]

	// The imported configurations must include the configuration from step 2.
	expected := map[string]string{
		"action":        "penalty",
		"spillsize":     "300",
		"bandwidth":     "400",
		"memsize":       "500",
		"broadcastsize": "600",
	}

	// For items not included in step 2 (configurations set in step 1 but reset in step 2), check if their values
	// ​​have been reset to -1.
	shouldMinusOne := map[string]bool{
		"blocktime":      true,
		"elapsedtime":    true,
		"allcputime":     true,
		"cpuskewpercent": true,
		"cpuavgpercent":  true,
	}

	cfgMap := make(map[string]string)
	for k, v := range is.Attributes {
		if strings.HasPrefix(k, "configurations.") && strings.HasSuffix(k, ".key") {
			parts := strings.Split(k, ".")
			if len(parts) < 3 {
				continue
			}

			valKey := fmt.Sprintf("configurations.%s.value", parts[1])
			cfgMap[v] = is.Attributes[valKey]
		}
	}
	for k, expectVal := range expected {
		actualVal, ok := cfgMap[k]
		if !ok {
			return fmt.Errorf("import check failed: missing key %s", k)
		}
		if actualVal != expectVal {
			return fmt.Errorf("import check failed: key %s expect %s, got %s", k, expectVal, actualVal)
		}
	}

	for k := range shouldMinusOne {
		actualVal, ok := cfgMap[k]
		if !ok {
			return fmt.Errorf("import check failed: missing old key %s (should be -1)", k)
		}
		if actualVal != "-1" {
			return fmt.Errorf("import check failed: key %s should be -1, got %s", k, actualVal)
		}
	}

	return nil
}

func testAccClusterExceptionRule_basic_step1(name string) string {
	return fmt.Sprintf(`
variable "exception_rule_configurations" {
  description = "The configurations of the exception rule."
  type        = list(object({
    key   = string
    value = string
  }))

  default = [
    {
      key   = "action"
      value = "abort"
    },
    {
      key   = "blocktime"
      value = "300"
    },
    {
      key   = "elapsedtime"
      value = "400"
    },
    {
      key   = "allcputime"
      value = "500"
    },
    {
      key   = "cpuskewpercent"
      value = "60"
    },
    {
      key   = "cpuavgpercent"
      value = "70"
    },
  ]
}

resource "huaweicloud_dws_cluster_exception_rule" "test" {
  cluster_id = "%[1]s"
  name       = "%[2]s"

  dynamic "configurations" {
    for_each = var.exception_rule_configurations

    content {
      key   = configurations.value.key
      value = configurations.value.value
    }
  }
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}

// Configuration values ​​configured in step 1 but not in step 2 will be reset to -1 (no limit).
func testAccClusterExceptionRule_basic_step2(name string) string {
	return fmt.Sprintf(`
variable "exception_rule_configurations" {
  description = "The configurations of the exception rule."
  type        = list(object({
    key   = string
    value = string
  }))

  default = [
    {
      key   = "action"
      value = "penalty"
    },
    {
      key   = "spillsize"
      value = "300"
    },
    {
      key   = "bandwidth"
      value = "400"
    },
    {
      key   = "memsize"
      value = "500"
    },
    {
      key   = "broadcastsize"
      value = "600"
    },
  ]
}

resource "huaweicloud_dws_cluster_exception_rule" "test" {
  cluster_id = "%[1]s"
  name       = "%[2]s"

  dynamic "configurations" {
    for_each = var.exception_rule_configurations

    content {
      key   = configurations.value.key
      value = configurations.value.value
    }
  }
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
