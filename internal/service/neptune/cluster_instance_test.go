package neptune_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/neptune"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfneptune "github.com/hashicorp/terraform-provider-aws/internal/service/neptune"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccNeptuneClusterInstance_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.cluster_instances"
	clusterResourceName := "aws_neptune_cluster.test"
	parameterGroupResourceName := "aws_neptune_parameter_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "rds", fmt.Sprintf("db:%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "true"),
					resource.TestMatchResourceAttr(resourceName, "availability_zone", regexp.MustCompile(fmt.Sprintf("^%s[a-z]{1}$", acctest.Region()))),
					resource.TestCheckResourceAttrPair(resourceName, "cluster_identifier", clusterResourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "dbi_resource_id"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttr(resourceName, "engine", "neptune"),
					resource.TestCheckResourceAttrSet(resourceName, "engine_version"),
					resource.TestCheckResourceAttr(resourceName, "identifier", rName),
					resource.TestCheckResourceAttrPair(resourceName, "instance_class", "data.aws_neptune_orderable_db_instance.test", "instance_class"),
					resource.TestCheckResourceAttr(resourceName, "kms_key_arn", ""),
					resource.TestCheckResourceAttrPair(resourceName, "neptune_parameter_group_name", parameterGroupResourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "neptune_subnet_group_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "port", strconv.Itoa(tfneptune.DefaultPort)),
					resource.TestCheckResourceAttrSet(resourceName, "preferred_backup_window"),
					resource.TestCheckResourceAttrSet(resourceName, "preferred_maintenance_window"),
					resource.TestCheckResourceAttr(resourceName, "promotion_tier", "3"),
					resource.TestCheckResourceAttr(resourceName, "publicly_accessible", "false"),
					resource.TestCheckResourceAttr(resourceName, "storage_encrypted", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "writer", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClusterInstanceConfig_modified(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "false"),
				),
			},
		},
	})
}

func TestAccNeptuneClusterInstance_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.cluster_instances"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfneptune.ResourceClusterInstance(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNeptuneClusterInstance_nameGenerated(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_nameGenerated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestMatchResourceAttr(resourceName, "identifier", regexp.MustCompile("^tf-")),
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

func TestAccNeptuneClusterInstance_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_namePrefix(rName, "tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "identifier", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "identifier_prefix", "tf-acc-test-prefix-"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"identifier_prefix",
				},
			},
		},
	})
}

func TestAccNeptuneClusterInstance_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.cluster_instances"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccClusterInstanceConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccClusterInstanceConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccNeptuneClusterInstance_withAZ(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.cluster_instances"
	availabiltyZonesDataSourceName := "data.aws_availability_zones.available"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_az(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone", availabiltyZonesDataSourceName, "names.0"),
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

func TestAccNeptuneClusterInstance_withSubnetGroup(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.test"
	subnetGroupResourceName := "aws_neptune_subnet_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_subnetGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "neptune_subnet_group_name", subnetGroupResourceName, "name"),
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

func TestAccNeptuneClusterInstance_kmsKey(t *testing.T) {
	ctx := acctest.Context(t)
	var v neptune.DBInstance
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_neptune_cluster_instance.cluster_instances"
	kmsKeyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, neptune.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterInstanceConfig_kmsKey(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_arn", kmsKeyResourceName, "arn"),
				),
			},
		},
	})
}

func testAccCheckClusterInstanceExists(ctx context.Context, n string, v *neptune.DBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Neptune Cluster Instance ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).NeptuneConn()

		output, err := tfneptune.FindClusterInstanceByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckClusterInstanceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).NeptuneConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_neptune_cluster_instance" {
				continue
			}

			_, err := tfneptune.FindClusterInstanceByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Neptune Cluster Instance %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccClusterInstanceConfig_baseSansCluster(rName string) string {
	return fmt.Sprintf(`
data "aws_neptune_orderable_db_instance" "test" {
  engine         = "neptune"
  engine_version = aws_neptune_cluster.test.engine_version
  license_model  = "amazon-license"

  preferred_instance_classes = ["db.t3.medium", "db.r5.large", "db.r4.large"]
}

resource "aws_neptune_parameter_group" "test" {
  name   = %[1]q
  family = "neptune1.2"

  parameter {
    name  = "neptune_query_timeout"
    value = "25"
  }
}
`, rName)
}

func testAccClusterInstanceConfig_base(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_baseSansCluster(rName), acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_neptune_cluster" "test" {
  cluster_identifier                   = %[1]q
  availability_zones                   = slice(data.aws_availability_zones.available.names, 0, min(3, length(data.aws_availability_zones.available.names)))
  engine                               = "neptune"
  neptune_cluster_parameter_group_name = "default.neptune1.2"
  skip_final_snapshot                  = true
}
`, rName))
}

func testAccClusterInstanceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_neptune_cluster_instance" "cluster_instances" {
  identifier                   = %[1]q
  cluster_identifier           = aws_neptune_cluster.test.id
  instance_class               = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version               = data.aws_neptune_orderable_db_instance.test.engine_version
  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
  promotion_tier               = "3"
}
`, rName))
}

func testAccClusterInstanceConfig_modified(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_neptune_cluster_instance" "cluster_instances" {
  identifier                   = %[1]q
  cluster_identifier           = aws_neptune_cluster.test.id
  instance_class               = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version               = data.aws_neptune_orderable_db_instance.test.engine_version
  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
  auto_minor_version_upgrade   = false
  promotion_tier               = "3"
}
`, rName))
}

func testAccClusterInstanceConfig_nameGenerated(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), `
resource "aws_neptune_cluster_instance" "test" {
  cluster_identifier = aws_neptune_cluster.test.id
  instance_class     = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version     = data.aws_neptune_orderable_db_instance.test.engine_version

  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
}
`)
}

func testAccClusterInstanceConfig_namePrefix(rName, prefix string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_neptune_cluster_instance" "test" {
  identifier_prefix  = %[1]q
  cluster_identifier = aws_neptune_cluster.test.id
  instance_class     = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version     = data.aws_neptune_orderable_db_instance.test.engine_version

  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
}
`, prefix))
}

func testAccClusterInstanceConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_neptune_cluster_instance" "cluster_instances" {
  identifier                   = %[1]q
  cluster_identifier           = aws_neptune_cluster.test.id
  instance_class               = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version               = data.aws_neptune_orderable_db_instance.test.engine_version
  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
  promotion_tier               = "3"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccClusterInstanceConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_neptune_cluster_instance" "cluster_instances" {
  identifier                   = %[1]q
  cluster_identifier           = aws_neptune_cluster.test.id
  instance_class               = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version               = data.aws_neptune_orderable_db_instance.test.engine_version
  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
  promotion_tier               = "3"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccClusterInstanceConfig_az(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_base(rName), fmt.Sprintf(`
resource "aws_neptune_cluster_instance" "cluster_instances" {
  identifier                   = %[1]q
  cluster_identifier           = aws_neptune_cluster.test.id
  instance_class               = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version               = data.aws_neptune_orderable_db_instance.test.engine_version
  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
  promotion_tier               = "3"
  availability_zone            = data.aws_availability_zones.available.names[0]
}
`, rName))
}

func testAccClusterInstanceConfig_subnetGroup(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_baseSansCluster(rName), acctest.ConfigVPCWithSubnets(rName, 2), fmt.Sprintf(`
resource "aws_neptune_cluster_instance" "test" {
  identifier         = %[1]q
  cluster_identifier = aws_neptune_cluster.test.id
  instance_class     = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version     = data.aws_neptune_orderable_db_instance.test.engine_version

  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
}

resource "aws_neptune_subnet_group" "test" {
  name       = %[1]q
  subnet_ids = aws_subnet.test[*].id
}

resource "aws_neptune_cluster" "test" {
  cluster_identifier                   = %[1]q
  neptune_subnet_group_name            = aws_neptune_subnet_group.test.name
  neptune_cluster_parameter_group_name = "default.neptune1.2"
  skip_final_snapshot                  = true
}
`, rName))
}

func testAccClusterInstanceConfig_kmsKey(rName string) string {
	return acctest.ConfigCompose(testAccClusterInstanceConfig_baseSansCluster(rName), acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description = %[1]q

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "kms-tf-1",
  "Statement": [
    {
      "Sid": "Enable IAM User Permissions",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": "kms:*",
      "Resource": "*"
    }
  ]
}
POLICY
}

resource "aws_neptune_cluster_instance" "cluster_instances" {
  identifier                   = %[1]q
  cluster_identifier           = aws_neptune_cluster.test.id
  instance_class               = data.aws_neptune_orderable_db_instance.test.instance_class
  engine_version               = data.aws_neptune_orderable_db_instance.test.engine_version
  neptune_parameter_group_name = aws_neptune_parameter_group.test.name
}

resource "aws_neptune_cluster" "test" {
  cluster_identifier  = %[1]q
  availability_zones  = slice(data.aws_availability_zones.available.names, 0, min(3, length(data.aws_availability_zones.available.names)))
  skip_final_snapshot = true
  storage_encrypted   = true
  kms_key_arn         = aws_kms_key.test.arn

  neptune_cluster_parameter_group_name = "default.neptune1.2"
}
`, rName))
}
