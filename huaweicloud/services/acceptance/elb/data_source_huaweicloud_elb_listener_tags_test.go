package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceListenerTags_basic(t *testing.T) {
	var (
		dataSourcename = "data.huaweicloud_elb_listener_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourcename)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceListenerTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourcename, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourcename, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourcename, "tags.0.values.#"),
				),
			},
		},
	})
}

const testAccDataSourceListenerTags_basic = `data "huaweicloud_elb_listener_tags" "test" {}`
