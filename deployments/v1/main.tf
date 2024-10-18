module "ebay" {
  source = "./ebay"

  allow_read_warehouse_policy_arn  = aws_iam_policy.allow_read_warehouse.arn
  allow_write_warehouse_policy_arn = aws_iam_policy.allow_write_warehouse.arn

  write_label_print_buffer_policy_arn = aws_iam_policy.allow_write_label_print_buffer.arn
}

module "web" {
  source = "./web"
}