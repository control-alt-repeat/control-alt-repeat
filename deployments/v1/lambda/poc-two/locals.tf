locals {
  function_name = "poc-two"
  archive_path  = "${path.module}/../../../../lambda-handler-${local.function_name}.zip"
}
