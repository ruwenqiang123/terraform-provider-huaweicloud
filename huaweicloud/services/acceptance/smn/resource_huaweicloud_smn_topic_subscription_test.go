package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTopicSubscription_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSmnSubscribeToken(t)
			acceptance.TestAccPrecheckSmnSubscribedTopicUrn(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTopicSubscription_basic(),
			},
		},
	})
}

func testTopicSubscription_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic_subscription" "test" {
  token     = "%[1]s"
  topic_urn = "%[2]s"
}
`, acceptance.HW_SMN_SUBSCRIBE_TOKEN, acceptance.HW_SMN_SUBSCRIBED_TOPIC_URN)
}
