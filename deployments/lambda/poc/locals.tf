locals {
  function_name = "poc"
  src_path      = "${path.module}/../cmd/aws/lambda/poc/lambda/${local.function_name}"

  binary_name  = local.function_name
  binary_path  = "${path.module}/../cmd/aws/lambda/poc/tf_generated/${local.binary_name}"
  archive_path = "${path.module}/../cmd/aws/lambda/poc/tf_generated/${local.function_name}.zip"
}
