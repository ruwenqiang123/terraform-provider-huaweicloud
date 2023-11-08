package deprecated

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/networks"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
				),
			},
			{
				ResourceName:      "huaweicloud_networking_subnet_v2.subnet_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingV2Port_noip(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_noip,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 1),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_timeout(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_createExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_createExtraDHCPOpts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				ResourceName:      "huaweicloud_networking_subnet_v2.subnet_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"fixed_ip",
				},
			},
		},
	})
}

func TestAccNetworkingV2Port_updateExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "1"),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("huaweicloud_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("huaweicloud_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("huaweicloud_networking_port_v2.port_1", &port),
					resource.TestCheckNoResourceAttr(
						"huaweicloud_networking_port_v2.port_1", "extra_dhcp_option"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_port_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortExists(n string, port *ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

func testAccCheckNetworkingV2PortCountFixedIPs(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.FixedIPs) != expected {
			return fmtp.Errorf("Expected %d Fixed IPs, got %d", expected, len(port.FixedIPs))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortCountAllowedAddressPairs(
	port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.AllowedAddressPairs) != expected {
			return fmtp.Errorf("Expected %d Allowed Address Pairs, got %d", expected, len(port.AllowedAddressPairs))
		}

		return nil
	}
}

const testAccNetworkingV2Port_basic = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2Port_noip = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  }
}
`

const testAccNetworkingV2Port_allowedAddressPairs = `
resource "huaweicloud_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "huaweicloud_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "huaweicloud_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${huaweicloud_networking_router_v2.vrrp_router.id}"
  subnet_id = "${huaweicloud_networking_subnet_v2.vrrp_subnet.id}"
}

resource "huaweicloud_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "huaweicloud_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "huaweicloud_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${huaweicloud_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${huaweicloud_networking_port_v2.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${huaweicloud_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${huaweicloud_networking_port_v2.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingV2Port_timeout = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Port_createExtraDHCPOpts = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionA"
    value = "valueA"
  }
  extra_dhcp_option {
    name = "optionB"
    value = "valueB"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_1 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_2 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }
  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_3 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }
  extra_dhcp_option {
    name = "optionE"
    value = "valueE"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_4 = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  fixed_ip {
    subnet_id =  "${huaweicloud_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`
