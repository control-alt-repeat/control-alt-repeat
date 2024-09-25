module "auth" {
    source = "./auth"
}

module "lambda_import_listing" {
    source = "./lambda/import-listing"
}
