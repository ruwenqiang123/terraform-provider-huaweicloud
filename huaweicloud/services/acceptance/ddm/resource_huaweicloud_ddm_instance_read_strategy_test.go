package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDdmInstanceReadStrategy_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()
	schemaName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ddm_instance_read_strategy.test"
	pwd := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDdmInstanceReadStrategy_basic(name, pwd, schemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_ddm_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "read_weights.#", "1"),
				),
			},
			{
				Config: testDdmInstanceReadStrategy_update(name, pwd, schemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_ddm_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "read_weights.#", "2"),
				),
			},
		},
	})
}

func testDdmInstanceReadStrategy_basic(name, pws, schemaName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_instance_read_strategy" "test" {
  depends_on  = [huaweicloud_ddm_schema.test]
  instance_id = huaweicloud_ddm_instance.test.id

  read_weights {
    db_id  = huaweicloud_rds_instance.test.id
    weight = 100
  }
}
`, testDdmInstanceReadStrategyBase(name, pws, schemaName))
}

func testDdmInstanceReadStrategy_update(name, pws, schemaName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ddm_instance_read_strategy" "test" {
  depends_on  = [huaweicloud_ddm_schema.test]
  instance_id = huaweicloud_ddm_instance.test.id

  read_weights {
    db_id  = huaweicloud_rds_instance.test.id
    weight = 60
  }

  read_weights {
    db_id  = huaweicloud_rds_read_replica_instance.test.id
    weight = 40
  }
}
`, testDdmInstanceReadStrategyBase(name, pws, schemaName))
}

func testDdmInstanceReadStrategyBase(name, pws, schemaName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_ddm_engines" test {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" test {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}
  
resource "huaweicloud_rds_instance" "test" {
  name              = "%[1]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  db {
    password = "%[2]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
  
data "huaweicloud_rds_flavors" "replica" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
  group_type    = "dedicated"
  memory        = 4
  vcpus         = 2
}
	
resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%[1]s-read-replica"
  flavor              = data.huaweicloud_rds_flavors.replica.flavors[0].name
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]

  volume {
    type              = "CLOUDSSD"
    size              = 50
    limit_size        = 400
    trigger_threshold = 10
  }
}

resource "huaweicloud_ddm_instance" "test" {
  depends_on        = [huaweicloud_rds_read_replica_instance.test]
  name              = "%[1]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_ddm_schema" "test" {
  instance_id  = huaweicloud_ddm_instance.test.id
  name         = "%[3]s"
  shard_mode   = "single"
  shard_number = "1"
  
  data_nodes {
    id             = huaweicloud_rds_instance.test.id
    admin_user     = "root"
    admin_password = "%[2]s"
  }
  
  delete_rds_data = "true"
  
  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}
`, name, pws, schemaName)
}
