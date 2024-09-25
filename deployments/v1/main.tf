module "lambda_poc" {
    source = "./lambda/poc"
}

module "lambda_poc_two" {
    source = "./lambda/poc-two"
}

module "lambda_ebay_import_listing" {
    source = "./lambda/ebay-import-listing"
}
