# resource "aws_lambda_function" "function" {
#   function_name    = "control-alt-repeat-v1-${local.function_name}"
#   description      = "It's a function to prove I am not repeating myself"
#   role             = aws_iam_role.lambda.arn
#   handler          = "bootstrap"
#   filename         = local.archive_path
#   # source_code_hash = filebase64sha256(local.archive_path)
#   runtime          = "provided.al2023"
# }

# resource "aws_cloudwatch_log_group" "log_group" {
#   name              = "/aws/lambda/${aws_lambda_function.function.function_name}"
#   retention_in_days = 7
# }
