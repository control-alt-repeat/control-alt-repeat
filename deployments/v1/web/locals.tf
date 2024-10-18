locals {
  function_name = "web"
  archive_path  = "${path.root}/../../lambda-handler-${local.function_name}.zip"
}
