# You cannot create a new backend by simply defining this and then
# immediately proceeding to "terraform apply". The S3 backend must
# be bootstrapped according to the simple yet essential procedure in
# https://github.com/cloudposse/terraform-aws-tfstate-backend#usage
module "terraform_state_backend" {
  source = "cloudposse/tfstate-backend/aws"
  # Cloud Posse recommends pinning every module to a specific version
  version    = "1.5.0"
  namespace  = "control-alt-repeat"
  stage      = "test"
  name       = "terraform"
  attributes = ["state"]

  terraform_backend_config_file_path = "."
  terraform_backend_config_file_name = "backend.tf"
  force_destroy                      = false
}

module "v1" {
  source = "./../../deployments/v1"
  stage_name = reverse(split("/", path.cwd))[0]
}

provider "aws" {
  region = "eu-west-2"

  default_tags {
    tags = {
      github_repository = "control-alt-repeat"
      github_organisation = "Control-Alt-Repeat"
      stage_name = "test"
    }
  }
}

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }

  required_version = ">= 1.9.6"
}
