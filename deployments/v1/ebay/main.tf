module "lambda_import_listing" {
  source = "./lambda/import-listing"

  read_listings_policy_arn  = aws_iam_policy.allow_read_listings.arn
  write_listings_policy_arn = aws_iam_policy.allow_write_listings.arn

  ebay_auth_ssm_access_policy_arn = aws_iam_policy.ebay_auth_ssm_access_policy.arn

  write_label_print_buffer_policy_arn = var.write_label_print_buffer_policy_arn
}

module "lambda_notification_endpoint" {
  source = "./lambda/notification-endpoint"

  allow_write_ebay_incoming_notifications_arn = allow_write_ebay_incoming_notifications.arn
}
