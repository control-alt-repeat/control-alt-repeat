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
