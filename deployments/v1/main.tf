module "lambda_poc" {
    source = "./lambda/poc"
}

module "lambda_poc_two" {
    source = "./lambda/poc-two"
}

module "ebay" {
    source = "./ebay"
}