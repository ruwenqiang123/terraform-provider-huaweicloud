package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/monitors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getMonitorResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud LB v2 client: %s", err)
	}
	resp, err := monitors.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the monitor (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "huaweicloud_lb_monitor.monitor_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
				),
			},
			{
				Config: testAccLBV2MonitorConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "3"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
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

func TestAccLBV2Monitor_udp(t *testing.T) {
	var monitor monitors.Monitor
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_monitor.monitor_udp"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_udp(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "UDP_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
				),
			},
		},
	})
}

func TestAccLBV2Monitor_http(t *testing.T) {
	var monitor monitors.Monitor
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_monitor.monitor_http"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&monitor,
		getMonitorResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_http(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/api"),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "expected_codes", "200-202"),
				),
			},
		},
	})
}

func testAccLBV2MonitorConfig_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}
`, rName, rName, rName)
}

func testAccLBV2MonitorConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_monitor" "monitor_1" {
  pool_id     = huaweicloud_lb_pool.pool_1.id
  name        = "%s"
  type        = "TCP"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
`, testAccLBV2MonitorConfig_base(rName), rName)
}

func testAccLBV2MonitorConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_monitor" "monitor_1" {
  pool_id     = huaweicloud_lb_pool.pool_1.id
  name        = "%s"
  type        = "TCP"
  delay       = 30
  timeout     = 15
  max_retries = 3
  port        = 8888
}
`, testAccLBV2MonitorConfig_base(rName), rNameUpdate)
}

func testAccLBV2MonitorConfig_http(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_monitor" "monitor_http" {
  pool_id        = huaweicloud_lb_pool.pool_1.id
  name           = "%s"
  type           = "HTTP"
  delay          = 20
  timeout        = 10
  max_retries    = 5
  url_path       = "/api"
  expected_codes = "200-202"
}
`, testAccLBV2MonitorConfig_base(rName), rName)
}

func testAccLBV2MonitorConfig_udp(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "UDP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "UDP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_lb_monitor" "monitor_udp" {
  pool_id     = huaweicloud_lb_pool.pool_1.id
  name        = "%s"
  type        = "UDP_CONNECT"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
`, rName, rName, rName, rName)
}
