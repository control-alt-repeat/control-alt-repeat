locals {
  function_name = "ebay-notification-endpoint"
  archive_path  = "${path.root}/../../lambda-handler-${local.function_name}.zip"
}
