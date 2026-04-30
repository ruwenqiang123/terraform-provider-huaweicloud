package obs

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OBS PUT ?publicAccessBlock
// @API OBS DELETE ?publicAccessBlock
// @API OBS GET ?publicAccessBlock
func ResourceObsBucketBpa() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketBpaCreate,
		UpdateContext: resourceObsBucketBpaUpdate,
		ReadContext:   resourceObsBucketBpaRead,
		DeleteContext: resourceObsBucketBpaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"bucket"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"block_public_acls": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ignore_public_acls": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"block_public_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"restrict_public_buckets": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceObsBucketBpaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		bucket = d.Get("bucket").(string)
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	opts := &obs.PutBucketPublicAccessBlockInput{
		Bucket: bucket,
		PublicAccessBlockConfiguration: obs.PublicAccessBlockConfiguration{
			BlockPublicAcls:       d.Get("block_public_acls").(bool),
			IgnorePublicAcls:      d.Get("ignore_public_acls").(bool),
			BlockPublicPolicy:     d.Get("block_public_policy").(bool),
			RestrictPublicBuckets: d.Get("restrict_public_buckets").(bool),
		},
	}

	_, err = obsClient.PutBucketPublicAccessBlock(opts)
	if err != nil {
		return diag.FromErr(getObsError("Error setting public access block of OBS bucket", bucket, err))
	}

	d.SetId(bucket)

	return resourceObsBucketBpaRead(ctx, d, meta)
}

func resourceObsBucketBpaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		bucket = d.Get("bucket").(string)
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	opts := &obs.PutBucketPublicAccessBlockInput{
		Bucket: bucket,
		PublicAccessBlockConfiguration: obs.PublicAccessBlockConfiguration{
			BlockPublicAcls:       d.Get("block_public_acls").(bool),
			IgnorePublicAcls:      d.Get("ignore_public_acls").(bool),
			BlockPublicPolicy:     d.Get("block_public_policy").(bool),
			RestrictPublicBuckets: d.Get("restrict_public_buckets").(bool),
		},
	}

	_, err = obsClient.PutBucketPublicAccessBlock(opts)
	if err != nil {
		return diag.FromErr(getObsError("Error updating public access block of OBS bucket", bucket, err))
	}

	return resourceObsBucketBpaRead(ctx, d, meta)
}

func resourceObsBucketBpaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	// Deleting BPA is equivalent to setting all fields to false, so this resource does not provide checkDelete `404`.
	output, err := obsClient.GetBucketPublicAccessBlock(d.Id())
	if err != nil {
		return diag.FromErr(getObsError("Error retrieving OBS bucket public access block", d.Id(), err))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bucket", d.Id()),
		d.Set("block_public_acls", output.BlockPublicAcls),
		d.Set("ignore_public_acls", output.IgnorePublicAcls),
		d.Set("block_public_policy", output.BlockPublicPolicy),
		d.Set("restrict_public_buckets", output.RestrictPublicBuckets),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OBS bucket public access block fields: %s", err)
	}

	return nil
}

func resourceObsBucketBpaDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	_, err = obsClient.DeleteBucketPublicAccessBlock(d.Id())
	if err != nil {
		return diag.FromErr(getObsError("Error deleting public access block of OBS bucket", d.Id(), err))
	}

	return nil
}
