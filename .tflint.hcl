# .tflint.hcl

plugin "terraform" {
    # Plugin common attributes
    enabled = true
    preset = "recommended"
}

plugin "aws" {
    enabled = true
    version = "0.34.0"
    source  = "github.com/terraform-linters/tflint-ruleset-aws"
}
