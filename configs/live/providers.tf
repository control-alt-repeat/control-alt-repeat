terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "5.73.0"
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

provider "aws" {
  region = "eu-west-2"

  default_tags {
    tags = {
      github_repository   = "control-alt-repeat"
      github_organisation = "control-alt-repeat"
      stage_name          = "test"
    }
  }
}
