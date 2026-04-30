package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getOBSBucketBpaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return nil, fmt.Errorf("error creating OBS Client: %s", err)
	}

	output, err := obsClient.GetBucketPublicAccessBlock(state.Primary.ID)
	if err != nil {
		return nil, err
	}

	blockPublicAcls := output.BlockPublicAcls
	ignorePublicAcls := output.IgnorePublicAcls
	blockPublicPolicy := output.BlockPublicPolicy
	restrictPublicBuckets := output.RestrictPublicBuckets
	if !blockPublicAcls && !ignorePublicAcls && !blockPublicPolicy && !restrictPublicBuckets {
		return nil, golangsdk.ErrDefault404{}
	}

	return output, nil
}

func TestAccObsBucketBpa_basic(t *testing.T) {
	var obj interface{}

	bucketName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_obs_bucket_bpa.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOBSBucketBpaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketBpa_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", bucketName),
					resource.TestCheckResourceAttr(rName, "block_public_acls", "true"),
					resource.TestCheckResourceAttr(rName, "ignore_public_acls", "true"),
					resource.TestCheckResourceAttr(rName, "block_public_policy", "true"),
					resource.TestCheckResourceAttr(rName, "restrict_public_buckets", "true"),
				),
			},
			{
				Config: testAccObsBucketBpa_update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", bucketName),
					resource.TestCheckResourceAttr(rName, "block_public_acls", "true"),
					resource.TestCheckResourceAttr(rName, "ignore_public_acls", "false"),
					resource.TestCheckResourceAttr(rName, "block_public_policy", "true"),
					resource.TestCheckResourceAttr(rName, "restrict_public_buckets", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccObsBucketBpa_base(bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%s"
  storage_class = "STANDARD"
  acl           = "private"
}
`, bucketName)
}

func testAccObsBucketBpa_basic(bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_bpa" "test" {
  bucket                  = huaweicloud_obs_bucket.test.bucket
  block_public_acls       = true
  ignore_public_acls      = true
  block_public_policy     = true
  restrict_public_buckets = true
}
`, testAccObsBucketBpa_base(bucketName))
}

func testAccObsBucketBpa_update(bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_bpa" "test" {
  bucket                  = huaweicloud_obs_bucket.test.bucket
  block_public_acls       = true
  ignore_public_acls      = false
  block_public_policy     = true
  restrict_public_buckets = false
}
`, testAccObsBucketBpa_base(bucketName))
}
