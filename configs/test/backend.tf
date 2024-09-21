terraform {
  required_version = ">= 1.0.0"

  backend "s3" {
    region  = "eu-west-2"
    bucket  = "control-alt-repeat-test-terraform-state"
    key     = "terraform.tfstate"
    profile = ""
    encrypt = "true"

    dynamodb_table = "control-alt-repeat-test-terraform-state-lock"
  }
}
