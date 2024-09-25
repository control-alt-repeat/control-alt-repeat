module "auth" {
    source = "./auth"
}

module "lambda_import_listing" {
    source = "./lambda/import-listing"

    read_listings_policy_arn = aws_iam_policy.function_read_listings.arn
    write_listings_policy_arn = aws_iam_policy.function_write_listings.arn

    read_tokens_policy_arn = module.auth.allow_read_tokens_policy_arn
    write_tokens_policy_arn = module.auth.allow_write_tokens_policy_arn
}
