package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/elb"
)

func getDomainResolutionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("elb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	return elb.GetDomainResolutionIpAddresses(client, state.Primary.ID)
}

func TestAccDomainResolution_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_elb_domain_resolution.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		name1        = fmt.Sprintf("test.pub%s.com", acctest.RandString(5))
		name2        = fmt.Sprintf("acc.int%s.com", acctest.RandString(5))
		zoneName1    = fmt.Sprintf("%s.", name1)
		zoneName2    = fmt.Sprintf("%s.", name2)
		obj          interface{}
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDomainResolutionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainResolutionConfig_basic(zoneName1, zoneName2, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id", "huaweicloud_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "public_domain_name_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "public_dns_zone_name", "huaweicloud_dns_zone.test1", "name"),
					resource.TestCheckResourceAttr(resourceName, "public_dns_record_set_ttl", "500"),
					resource.TestCheckResourceAttr(resourceName, "private_domain_name_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "private_dns_zone_name", "huaweicloud_dns_zone.test2", "name"),
					resource.TestCheckResourceAttr(resourceName, "private_dns_zone_type", "private"),
					resource.TestCheckResourceAttr(resourceName, "private_dns_record_set_ttl", "400"),
				),
			},
			{
				Config: testAccDomainResolutionConfig_update(zoneName1, zoneName2, rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "public_domain_name_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "private_domain_name_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "private_dns_zone_name", "huaweicloud_dns_zone.test1", "name"),
					resource.TestCheckResourceAttr(resourceName, "private_dns_zone_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "private_dns_record_set_ttl", "200"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"loadbalancer_id",
					"public_domain_name_enable",
					"public_dns_zone_name",
					"public_dns_record_set_ttl",
					"private_domain_name_enable",
					"private_dns_zone_name",
					"private_dns_zone_type",
					"private_dns_record_set_ttl",
				},
			},
		},
	})
}

func testAccDomainResolutionConfig_base(zoneName1, zoneName2, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_zone" "test1" {
  name      = "%[2]s"
  zone_type = "public"
}

resource "huaweicloud_dns_zone" "test2" {
  name      = "%[3]s"
  zone_type = "private"

  router {
    router_id = huaweicloud_vpc.test.id
  }
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[4]s"
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5

  lifecycle {
    ignore_changes = [
      l4_flavor_id, l7_flavor_id
    ]
  }
}
`, common.TestVpc(name), zoneName1, zoneName2, name)
}

func testAccDomainResolutionConfig_basic(zoneName1, zoneName2, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_domain_resolution" "test" {
  loadbalancer_id            = huaweicloud_elb_loadbalancer.test.id
  public_domain_name_enable  = true
  public_dns_zone_name       = huaweicloud_dns_zone.test1.name
  public_dns_record_set_ttl  = 500
  private_domain_name_enable = true
  private_dns_zone_name      = huaweicloud_dns_zone.test2.name
  private_dns_zone_type      = "private"
  private_dns_record_set_ttl = 400
}
`, testAccDomainResolutionConfig_base(zoneName1, zoneName2, name))
}

func testAccDomainResolutionConfig_update(zoneName1, zoneName2, name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_domain_resolution" "test" {
  loadbalancer_id            = huaweicloud_elb_loadbalancer.test.id
  public_domain_name_enable  = false
  private_domain_name_enable = true
  private_dns_zone_name      = huaweicloud_dns_zone.test1.name
  private_dns_zone_type      = "public"
  private_dns_record_set_ttl = 200
}
`, testAccDomainResolutionConfig_base(zoneName1, zoneName2, name))
}
