locals {
  function_name = "ebay-import-listing"
  archive_path = "${path.module}/../../../../lambda-handler-${local.function_name}.zip"
}
