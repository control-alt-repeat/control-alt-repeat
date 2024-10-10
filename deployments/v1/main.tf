module "lambda_poc" {
  source = "./lambda/poc"
}

module "lambda_poc_two" {
  source = "./lambda/poc-two"
}

module "ebay" {
  source = "./ebay"

  write_label_print_buffer_policy_arn = aws_iam_policy.allow_write_label_print_buffer.arn
}