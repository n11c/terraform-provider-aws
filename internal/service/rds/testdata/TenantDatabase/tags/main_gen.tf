ephemeral "aws_secretsmanager_random_password" "test" {
  password_length     = 20
  exclude_punctuation = true
}

data "aws_rds_engine_version" "test" {
  engine = "oracle-se2-cdb"
}

data "aws_rds_orderable_db_instance" "test" {
  engine         = data.aws_rds_engine_version.test.engine
  engine_version = data.aws_rds_engine_version.test.version
  license_model  = "license-included"
  storage_type   = "gp2"

  preferred_instance_classes = ["db.m5.xlarge", "db.m6i.xlarge", "db.m5.2xlarge", "db.m6i.2xlarge"]
}

resource "aws_db_instance" "test" {
  identifier          = var.rName
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

resource "aws_db_tenant_database" "test" {
  db_instance_identifier = aws_db_instance.test.identifier
  tenant_db_name         = "TESTPDB"
  master_username        = "tfadmin"
  master_user_password   = ephemeral.aws_secretsmanager_random_password.test.random_password
  skip_final_snapshot    = true

  tags = var.resource_tags
}

variable "rName" {
  description = "Name for resource"
  type        = string
  nullable    = false
}

variable "resource_tags" {
  description = "Tags to set on resource. To specify no tags, set to `null`"
  # Not setting a default, so that this must explicitly be set to `null` to specify no tags
  type     = map(string)
  nullable = true
}
