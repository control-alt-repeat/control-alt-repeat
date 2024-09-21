locals {
  function_name = "poc"
  src_path      = "${path.module}/../../../cmd/aws/lambda/${local.function_name}"

  binary_name  = local.function_name
  binary_path  = "${local.src_path}/tf_generated/${local.binary_name}"
  archive_path = "${local.src_path}/tf_generated/${local.function_name}.zip"
}
