resource "aws_lambda_function" "function" {
  function_name    = "control-alt-repeat-v1-${local.function_name}"
  description      = "It's a function to prove deployment capability"
  role             = aws_iam_role.lambda.arn
  handler          = "bootstrap"
  filename         = local.archive_path
  source_code_hash = filebase64sha256(local.archive_path)
  runtime          = "provided.al2023"
  timeout          = 60
}

resource "aws_lambda_function_url" "latest" {
  function_name      = aws_lambda_function.function.function_name
  authorization_type = "NONE"
}

resource "aws_ssm_parameter" "web_endpoint_parameter" {
  name  = "/control_alt_repeat/ebay/live/web/endpoint"
  type  = "String"
  value = aws_lambda_function_url.latest.function_url
}

resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 7
}
