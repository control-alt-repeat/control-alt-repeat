locals {
  function_name = "poc"
  src_path      = "${path.module}/../../../cmd/aws/lambda/${local.function_name}"

  binary_name  = local.function_name
  binary_path  = "${local.src_path}/${local.binary_name}/bootstrap"
  archive_path = "${local.src_path}/bootstrap.zip"
}
