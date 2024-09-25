output "allow_read_tokens_policy_arn" {
  value = aws_iam_policy.function_read_tokens.arn
}

output "allow_write_tokens_policy_arn" {
  value = aws_iam_policy.function_write_tokens.arn
}
