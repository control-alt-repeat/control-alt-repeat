output "allow_read_listings_arn" {
    value = aws_iam_policy.allow_read_listings.arn
}

output "allow_write_listings_arn" {
    value = aws_iam_policy.allow_write_listings.arn
}

output "ebay_auth_ssm_access_policy_arn" {
    value = aws_iam_policy.ebay_auth_ssm_access_policy.arn
}