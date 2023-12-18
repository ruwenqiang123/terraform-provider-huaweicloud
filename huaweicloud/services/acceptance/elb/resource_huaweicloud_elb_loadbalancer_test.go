package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getELBResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}

	eipID := state.Primary.Attributes["ipv4_eip_id"]
	eipType := state.Primary.Attributes["iptype"]
	if eipType != "" && eipID != "" {
		eipClient, err := c.NetworkingV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
		}

		if _, err := eips.Get(eipClient, eipID).Extract(); err != nil {
			return nil, err
		}
	}

	return loadbalancers.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "backend_subnets.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "backend_subnets.0",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttr(resourceName, "backend_subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEpsId(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEIP(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_withEIP_Bandwidth_Id(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_withEIP_Bandwidth_Id(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "iptype", "5_bgp"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv4_eip_id"),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_prePaid(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_prePaidUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", "update flavors"),
					resource.TestCheckResourceAttrPair(
						resourceName, "l4_flavor_id", "data.huaweicloud_elb_flavors.l4flavors", "ids.0"),
					resource.TestCheckResourceAttrPair(
						resourceName, "l7_flavor_id", "data.huaweicloud_elb_flavors.l7flavors", "ids.0"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.huaweicloud_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccElbV3LoadBalancer_availabilityZone(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getELBResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3LoadBalancerConfig_availabilityZone(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3LoadBalancerConfig_availabilityZoneUpdate(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.huaweicloud_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccElbV3LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_subnet" "test_1" {
  name       = "%[2]s_1"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
	
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  protection_status = "nonProtection"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3LoadBalancerConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_subnet" "test_1" {
  name       = "%[2]s_1"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[3]s"
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id,
    huaweicloud_vpc_subnet.test_1.id,
  ]

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, common.TestVpc(rName), rName, rNameUpdate)
}

func testAccElbV3LoadBalancerConfig_withEpsId(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name                  = "%s"
  ipv4_subnet_id        = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
  enterprise_project_id = "%s"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccElbV3LoadBalancerConfig_withEIP(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%s"
  ipv4_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  iptype                = "5_bgp"
  bandwidth_charge_mode = "traffic"
  sharetype             = "PER"
  bandwidth_size        = 5
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_withEIP_Bandwidth_Id(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[1]s"
  ipv4_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  iptype       = "5_bgp"
  sharetype    = "WHOLE"
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_prePaid(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = data.huaweicloud_vpc_subnet.test.id

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_prePaidUpdate(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_elb_flavors" "l4flavors" {
  type            = "L4"
  max_connections = 1000000
  cps             = 20000
  bandwidth       = 100
}

data "huaweicloud_elb_flavors" "l7flavors" {
  type            = "L7"
  max_connections = 400000
  cps             = 4000
  bandwidth       = 100
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = data.huaweicloud_vpc_subnet.test.id
  description     = "update flavors"
  l4_flavor_id    = data.huaweicloud_elb_flavors.l4flavors.ids[0]
  l7_flavor_id    = data.huaweicloud_elb_flavors.l7flavors.ids[0]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[1]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3LoadBalancerConfig_availabilityZone(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
	
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3LoadBalancerConfig_availabilityZoneUpdate(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[1]
  ]

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, common.TestVpc(rName), rNameUpdate)
}
