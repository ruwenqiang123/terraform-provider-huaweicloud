package cfw

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getProtectionRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getProtectionRule: Query the CFW Protection Rule detail
	var (
		getProtectionRuleHttpUrl = "v1/{project_id}/acl-rules"
		getProtectionRuleProduct = "cfw"
	)
	getProtectionRuleClient, err := conf.NewServiceClient(getProtectionRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ProtectionRule Client: %s", err)
	}

	getProtectionRulePath := getProtectionRuleClient.Endpoint + getProtectionRuleHttpUrl
	getProtectionRulePath = strings.ReplaceAll(getProtectionRulePath, "{project_id}", getProtectionRuleClient.ProjectID)

	getProtectionRulequeryParams := buildGetProtectionRuleQueryParams(state)
	getProtectionRulePath += getProtectionRulequeryParams

	getPotectionRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getProtectionRuleResp, err := getProtectionRuleClient.Request("GET", getProtectionRulePath, &getPotectionRulesOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving protection rule: %s", err)
	}

	getProtectionRuleRespBody, err := utils.FlattenResponse(getProtectionRuleResp)
	if err != nil {
		return nil, err
	}

	rules, err := jmespath.Search("data.records", getProtectionRuleRespBody)
	if err != nil {
		diag.Errorf("error parsing data.records from response= %#v", getProtectionRuleRespBody)
	}

	return cfw.FilterRules(rules.([]interface{}), state.Primary.ID)
}

func buildGetProtectionRuleQueryParams(state *terraform.ResourceState) string {
	res := "?offset=0&limit=1024"
	res = fmt.Sprintf("%s&object_id=%v", res, state.Primary.Attributes["object_id"])

	return res
}

func TestAccProtectionRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_protection_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectionRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProtectionRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "type", "0"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
					resource.TestCheckResourceAttr(rName, "action_type", "0"),
					resource.TestCheckResourceAttr(rName, "long_connect_enable", "0"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.address", "1.1.1.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.address", "1.1.1.2"),
					resource.TestCheckResourceAttrSet(rName, "rule_hit_count"),
				),
			},
			{
				Config: testProtectionRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(rName, "action_type", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.address", "2.2.2.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.address", "2.2.2.2"),
				),
			},
			{
				Config: testProtectionRule_region_list(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "direction", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.address", "2.2.2.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.type", "3"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.#", "3"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.0.description_en", "Greece"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.0.description_cn", "希腊"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.0.region_id", "GR"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.0.region_type", "0"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.1.description_en", "ZHEJIANG"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.1.description_cn", "浙江"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.1.region_id", "ZJ"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.1.region_type", "1"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.2.description_en", "Africa"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.2.description_cn", "非洲"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.2.region_id", "AF"),
					resource.TestCheckResourceAttr(rName, "destination.0.region_list.2.region_type", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testProtectionRuleImportState(rName),
				ImportStateVerifyIgnore: []string{
					"sequence", "type",
				},
			},
		},
	})
}

func TestAccProtectionRule_withTag_domainSet_customService(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_protection_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectionRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccProtectionRule_withTag_domainSet_customService_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "type", "0"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
					resource.TestCheckResourceAttr(rName, "action_type", "0"),
					resource.TestCheckResourceAttr(rName, "long_connect_enable", "0"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "direction", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.address", "1.1.1.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.domain_set_name", name+"_dg1"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.#", "2"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.0.source_port", "80"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.0.dest_port", "80"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.1.source_port", "8080"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.1.dest_port", "8080"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "rule_hit_count"),
				),
			},
			{
				Config: testAccProtectionRule_withTag_domainSet_customService_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.#", "3"),
					resource.TestCheckResourceAttr(rName, "destination.0.domain_set_name", name+"_dg2"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.0.source_port", "80"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.0.dest_port", "80"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.1.source_port", "8080"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.1.dest_port", "8080"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.2.source_port", "443"),
					resource.TestCheckResourceAttr(rName, "service.0.custom_service.2.dest_port", "443"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testProtectionRuleImportState(rName),
				ImportStateVerifyIgnore: []string{
					"sequence", "type",
				},
			},
		},
	})
}

func TestAccProtectionRule_withIpAddress_addressGroup_serviceGroup(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_cfw_protection_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProtectionRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttr(rName, "type", "0"),
					resource.TestCheckResourceAttr(rName, "address_type", "0"),
					resource.TestCheckResourceAttr(rName, "action_type", "0"),
					resource.TestCheckResourceAttr(rName, "long_connect_enable", "0"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.ip_address.#", "2"),
					resource.TestCheckResourceAttr(rName, "source.0.ip_address.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(rName, "source.0.ip_address.1", "1.1.1.2"),
					resource.TestCheckResourceAttr(rName, "destination.0.address_group.#", "2"),
					resource.TestCheckResourceAttr(rName, "service.0.service_group.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "rule_hit_count"),
				),
			},
			{
				Config: testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source.0.ip_address.#", "3"),
					resource.TestCheckResourceAttr(rName, "source.0.ip_address.0", "1.1.1.3"),
					resource.TestCheckResourceAttr(rName, "source.0.ip_address.1", "1.1.1.4"),
					resource.TestCheckResourceAttr(rName, "source.0.ip_address.2", "1.1.1.6"),
					resource.TestCheckResourceAttr(rName, "destination.0.address_group.#", "2"),
					resource.TestCheckResourceAttr(rName, "service.0.service_group.#", "3"),
				),
			},
			{
				Config: testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source.0.address_group.#", "2"),
					resource.TestCheckResourceAttr(rName, "destination.0.ip_address.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.ip_address.1", "1.1.1.2"),
					resource.TestCheckResourceAttr(rName, "service.0.protocol", "6"),
					resource.TestCheckResourceAttr(rName, "service.0.source_port", "8001"),
					resource.TestCheckResourceAttr(rName, "service.0.dest_port", "8002"),
				),
			},
			{
				Config: testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step4(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source.0.address_group.#", "2"),
					resource.TestCheckResourceAttr(rName, "destination.0.ip_address.0", "2.1.1.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.ip_address.1", "2.1.1.2"),
					resource.TestCheckResourceAttr(rName, "service.0.protocol", "6"),
					resource.TestCheckResourceAttr(rName, "service.0.source_port", "8001"),
					resource.TestCheckResourceAttr(rName, "service.0.dest_port", "8002"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testProtectionRuleImportState(rName),
				ImportStateVerifyIgnore: []string{
					"sequence", "type",
				},
			},
		},
	})
}

func testProtectionRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["object_id"] == "" {
			return "", fmt.Errorf("Attribute (object_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("Attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["object_id"] + "/" +
			rs.Primary.ID, nil
	}
}

func testProtectionRule_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type    = 0
    address = "1.1.1.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testProtectionRule_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%s-update"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test update"
  type                = 0
  address_type        = 0
  action_type         = 1
  long_connect_enable = 0
  status              = 1

  source {
    type    = 0
    address = "2.2.2.1"
  }

  destination {
    type    = 0
    address = "2.2.2.2"
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testProtectionRule_region_list(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[2]s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test update"
  type                = 0
  address_type        = 0
  action_type         = 1
  long_connect_enable = 0
  status              = 1
  direction           = 1

  source {
    type    = 0
    address = "2.2.2.1"
  }

  destination {
    type = 3

    region_list {
      description_cn = "希腊"
      description_en = "Greece"
      region_id      = "GR"
      region_type    = 0
    }

    region_list {
      description_cn = "浙江"
      description_en = "ZHEJIANG"
      region_id      = "ZJ"
      region_type    = 1
    }

    region_list {
      description_cn = "非洲"
      description_en = "Africa"
      region_id      = "AF"
      region_type    = 2
    }
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), name)
}

func testAccProtectionRule_withTag_domainSet_customService_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1
  direction           = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type            = 4
    domain_set_id   = huaweicloud_cfw_domain_name_group.dg1.id
    domain_set_name = huaweicloud_cfw_domain_name_group.dg1.name
  }

  service {
    type = 2

    custom_service {
      protocol    = 6
      source_port = 80
      dest_port   = 80			
    }

    custom_service {
      protocol    = 6
      source_port = 8080
      dest_port   = 8080
    }
  }

  sequence {
    top = 1
  }

  tags = {
    key = "value"
  }
}
`, testAccDatasourceFirewalls_basic(), testAccProtectionRule_advanced_base(name), name)
}

func testAccProtectionRule_withTag_domainSet_customService_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1
  direction           = 1

  source {
    type    = 0
    address = "1.1.1.1"
  }

  destination {
    type            = 4
    domain_set_id   = huaweicloud_cfw_domain_name_group.dg2.id
    domain_set_name = huaweicloud_cfw_domain_name_group.dg2.name
  }

  service {
    type = 2

    custom_service {
      protocol    = 6
      source_port = 80
      dest_port   = 80			
    }

    custom_service {
      protocol    = 6
      source_port = 8080
      dest_port   = 8080
    }

    custom_service {
      protocol    = 6
      source_port = 443
      dest_port   = 443
    }
  }

  sequence {
    top = 1
  }

  tags = {
    foo = "bar"
  }
}
`, testAccDatasourceFirewalls_basic(), testAccProtectionRule_advanced_base(name), name)
}

func testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type       = 5
    ip_address = ["1.1.1.1","1.1.1.2"]
  }

  destination {
    type = 5

    address_group = [
      huaweicloud_cfw_address_group.g1.id,
      huaweicloud_cfw_address_group.g2.id,
    ]
  }

  service {
    type = 2

    service_group = [
      huaweicloud_cfw_service_group.s1.id,
      huaweicloud_cfw_service_group.s2.id,
    ]
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), testAccProtectionRule_advanced_base(name), name)
}

func testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type       = 5
    ip_address = ["1.1.1.3","1.1.1.4","1.1.1.6"]
  }

  destination {
    type = 5

    address_group = [
      huaweicloud_cfw_address_group.g2.id,
      huaweicloud_cfw_address_group.g3.id,
    ]
  }

  service {
    type = 2

    service_group = [
      huaweicloud_cfw_service_group.s1.id,
      huaweicloud_cfw_service_group.s2.id,
      huaweicloud_cfw_service_group.s3.id,
    ]
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), testAccProtectionRule_advanced_base(name), name)
}

func testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type = 5

    address_group = [
      huaweicloud_cfw_address_group.g2.id,
      huaweicloud_cfw_address_group.g3.id,
    ]
  }

  destination {
    type       = 5
    ip_address = ["1.1.1.1","1.1.1.2"]
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), testAccProtectionRule_advanced_base(name), name)
}

func testAccProtectionRule_withIpAddress_addressGroup_serviceGroup_step4(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cfw_protection_rule" "test" {
  name                = "%[3]s"
  object_id           = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  description         = "terraform test"
  type                = 0
  address_type        = 0
  action_type         = 0
  long_connect_enable = 0
  status              = 1

  source {
    type = 5

    address_group = [
      huaweicloud_cfw_address_group.g1.id,
      huaweicloud_cfw_address_group.g2.id,
    ]
  }

  destination {
    type       = 5
    ip_address = ["2.1.1.1","2.1.1.2"]
  }

  service {
    type        = 0
    protocol    = 6
    source_port = 8001
    dest_port   = 8002
  }

  sequence {
    top = 1
  }
}
`, testAccDatasourceFirewalls_basic(), testAccProtectionRule_advanced_base(name), name)
}

func testAccProtectionRule_advanced_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_address_group" "g1" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[1]s_ag1"
  description = "address group 1"
}

resource "huaweicloud_cfw_address_group" "g2" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[1]s_ag2"
  description = "address group 2"
}

resource "huaweicloud_cfw_address_group" "g3" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[1]s_ag3"
  description = "address group 3"
}

resource "huaweicloud_cfw_service_group" "s1" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[1]s_sg1"
  description = "service group 1"
}

resource "huaweicloud_cfw_service_group" "s2" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[1]s_sg2"
  description = "service group 2"
}

resource "huaweicloud_cfw_service_group" "s3" {
  object_id   = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name        = "%[1]s_sg3"
  description = "service group 3"
}

resource "huaweicloud_cfw_domain_name_group" "dg1" {
  fw_instance_id = data.huaweicloud_cfw_firewalls.test.records[0].fw_instance_id
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[1]s_dg1"
  type           = 0
  description    = "created by terraform"
  
  domain_names {
    domain_name = "www.cfw-test1.com"
    description = "test domain 1"
  }
}

resource "huaweicloud_cfw_domain_name_group" "dg2" {
  fw_instance_id = data.huaweicloud_cfw_firewalls.test.records[0].fw_instance_id
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[1]s_dg2"
  type           = 0
  description    = "created by terraform"
  
  domain_names {
    domain_name = "www.cfw-test2.com"
    description = "test domain 2"
  }
}
`, name)
}
