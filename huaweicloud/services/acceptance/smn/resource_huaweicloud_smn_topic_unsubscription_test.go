package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTopicUnsubscription_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSmnSubscribedTopicUrn(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTopicUnsubscription_basic(),
			},
		},
	})
}

func testTopicUnsubscription_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic_unsubscription" "test" {
  subscription_urn = "%s"
}
`, acceptance.HW_SMN_SUBSCRIBED_TOPIC_URN)
}
