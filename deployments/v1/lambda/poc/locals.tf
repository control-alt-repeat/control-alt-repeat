locals {
  function_name = "poc"
  archive_path = "${path.module}/../../../../lambda-handler-${local.function_name}.zip"
}
