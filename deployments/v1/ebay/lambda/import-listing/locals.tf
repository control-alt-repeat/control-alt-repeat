locals {
  function_name = "ebay-import-listing"
  archive_path = "${path.root}/../../lambda-handler-${local.function_name}.zip"
}
