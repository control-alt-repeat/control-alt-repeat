module "lambda_import_listing" {
    source = "./lambda/import-listing"

    read_listings_policy_arn = aws_iam_policy.allow_read_listings.arn
    write_listings_policy_arn = aws_iam_policy.allow_write_listings.arn

    ebay_auth_ssm_access_policy_arn = aws_iam_policy.ebay_auth_ssm_access_policy.arn
}
