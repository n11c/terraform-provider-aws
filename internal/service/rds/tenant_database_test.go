package rds_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/retry"
	tfrds "github.com/hashicorp/terraform-provider-aws/internal/service/rds"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccRDSTenantDatabase_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := acctest.RandomWithPrefix(t, acctest.ResourcePrefix)
	resourceName := "aws_db_tenant_database.test"

	acctest.ParallelTest(ctx, t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RDSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTenantDatabaseDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTenantDatabaseConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTenantDatabaseExists(ctx, resourceName),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrARN),
					resource.TestCheckResourceAttr(resourceName, "db_instance_identifier", rName),
					resource.TestCheckResourceAttr(resourceName, "tenant_db_name", "TESTPDB"),
					resource.TestCheckResourceAttr(resourceName, "master_username", "tfadmin"),
					resource.TestCheckResourceAttrSet(resourceName, "dbi_resource_id"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_database_resource_id"),
					resource.TestCheckResourceAttrSet(resourceName, "tenant_database_create_time"),
					resource.TestCheckResourceAttr(resourceName, "skip_final_snapshot", acctest.CtTrue),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTenantDatabaseImportStateIdFunc(resourceName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"master_user_password",
					"skip_final_snapshot",
					"final_db_snapshot_identifier",
				},
			},
		},
	})
}

func TestAccRDSTenantDatabase_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := acctest.RandomWithPrefix(t, acctest.ResourcePrefix)
	resourceName := "aws_db_tenant_database.test"

	acctest.ParallelTest(ctx, t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RDSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTenantDatabaseDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTenantDatabaseConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantDatabaseExists(ctx, resourceName),
					acctest.CheckFrameworkResourceDisappears(ctx, t, tfrds.ResourceTenantDatabase, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRDSTenantDatabase_update(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := acctest.RandomWithPrefix(t, acctest.ResourcePrefix)
	resourceName := "aws_db_tenant_database.test"

	acctest.ParallelTest(ctx, t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RDSServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTenantDatabaseDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTenantDatabaseConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantDatabaseExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tenant_db_name", "TESTPDB"),
				),
			},
			{
				Config: testAccTenantDatabaseConfig_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTenantDatabaseExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tenant_db_name", "NEWPDB"),
				),
			},
		},
	})
}

func testAccCheckTenantDatabaseDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RDSClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_db_tenant_database" {
				continue
			}

			instanceID := rs.Primary.Attributes["db_instance_identifier"]
			tenantDBName := rs.Primary.Attributes["tenant_db_name"]

			_, err := tfrds.FindTenantDatabaseByName(ctx, conn, instanceID, tenantDBName)

			if retry.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("RDS Tenant Database %s/%s still exists", instanceID, tenantDBName)
		}

		return nil
	}
}

func testAccCheckTenantDatabaseExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RDSClient(ctx)

		instanceID := rs.Primary.Attributes["db_instance_identifier"]
		tenantDBName := rs.Primary.Attributes["tenant_db_name"]

		_, err := tfrds.FindTenantDatabaseByName(ctx, conn, instanceID, tenantDBName)

		return err
	}
}

func testAccTenantDatabaseImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		return fmt.Sprintf("%s/%s",
			rs.Primary.Attributes["db_instance_identifier"],
			rs.Primary.Attributes["tenant_db_name"],
		), nil
	}
}

// testAccTenantDatabaseConfig_base creates an Oracle SE2-CDB multi-tenant instance.
func testAccTenantDatabaseConfig_base(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigRandomPassword(),
		fmt.Sprintf(`
data "aws_rds_engine_version" "test" {
  engine = %[1]q
}

data "aws_rds_orderable_db_instance" "test" {
  engine         = data.aws_rds_engine_version.test.engine
  engine_version = data.aws_rds_engine_version.test.version
  license_model  = "license-included"
  storage_type   = "gp2"

  preferred_instance_classes = [%[2]s]
}

resource "aws_db_instance" "test" {
  identifier          = %[3]q
  engine              = data.aws_rds_orderable_db_instance.test.engine
  engine_version      = data.aws_rds_orderable_db_instance.test.engine_version
  instance_class      = data.aws_rds_orderable_db_instance.test.instance_class
  allocated_storage   = 20
  storage_type        = "gp2"
  license_model       = "license-included"
  multi_tenant        = true
  username            = "tfadmin"
  password_wo         = ephemeral.aws_secretsmanager_random_password.test.random_password
  password_wo_version = 1
  skip_final_snapshot = true
}
`, tfrds.InstanceEngineOracleStandard2CDB, strings.Replace(mainInstanceClasses, "db.t3.small", "frodo", 1), rName))
}

func testAccTenantDatabaseConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccTenantDatabaseConfig_base(rName),
		`
resource "aws_db_tenant_database" "test" {
  db_instance_identifier = aws_db_instance.test.identifier
  tenant_db_name         = "TESTPDB"
  master_username        = "tfadmin"
  master_user_password   = ephemeral.aws_secretsmanager_random_password.test.random_password
  skip_final_snapshot    = true
}
`)
}

func testAccTenantDatabaseConfig_updated(rName string) string {
	return acctest.ConfigCompose(
		testAccTenantDatabaseConfig_base(rName),
		`
resource "aws_db_tenant_database" "test" {
  db_instance_identifier = aws_db_instance.test.identifier
  tenant_db_name         = "NEWPDB"
  master_username        = "tfadmin"
  master_user_password   = ephemeral.aws_secretsmanager_random_password.test.random_password
  skip_final_snapshot    = true
}
`)
}
