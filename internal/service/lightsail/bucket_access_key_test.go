package lightsail_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/lightsail"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tflightsail "github.com/hashicorp/terraform-provider-aws/internal/service/lightsail"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccLightsailBucketAccessKey_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lightsail_bucket_access_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, lightsail.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, lightsail.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketAccessKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketAccessKeyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketAccessKeyExists(ctx, resourceName),
					resource.TestMatchResourceAttr(resourceName, "access_key_id", regexp.MustCompile(`((?:ASIA|AKIA|AROA|AIDA)([A-Z0-7]{16}))`)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestMatchResourceAttr(resourceName, "secret_access_key", regexp.MustCompile(`([a-zA-Z0-9+/]{40})`)),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_access_key", "bucket_name"},
			},
		},
	})
}

func TestAccLightsailBucketAccessKey_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_lightsail_bucket_access_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, lightsail.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, lightsail.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketAccessKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketAccessKeyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketAccessKeyExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tflightsail.ResourceBucketAccessKey(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckBucketAccessKeyExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LightsailConn()

		out, err := tflightsail.FindBucketAccessKeyById(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		if out == nil {
			return fmt.Errorf("BucketAccessKey %q does not exist", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckBucketAccessKeyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).LightsailConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_lightsail_bucket_access_key" {
				continue
			}

			_, err := tflightsail.FindBucketAccessKeyById(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return create.Error(names.Lightsail, create.ErrActionCheckingDestroyed, tflightsail.ResBucketAccessKey, rs.Primary.ID, errors.New("still exists"))
		}

		return nil
	}
}

func testAccBucketAccessKeyConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_lightsail_bucket" "test" {
  name      = %[1]q
  bundle_id = "small_1_0"
}
`, rName)
}

func testAccBucketAccessKeyConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccBucketAccessKeyConfig_base(rName), `
resource "aws_lightsail_bucket_access_key" "test" {
  bucket_name = aws_lightsail_bucket.test.id
}
`)
}
