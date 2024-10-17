resource "aws_lambda_function" "function" {
  function_name    = "control-alt-repeat-v1-${local.function_name}"
  description      = "It's a function to prove deployment capability"
  role             = aws_iam_role.lambda.arn
  handler          = "bootstrap"
  filename         = local.archive_path
  source_code_hash = filebase64sha256(local.archive_path)
  runtime          = "provided.al2023"
  timeout          = 10
}

resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 7
}

resource "aws_lambda_permission" "allow_notification_bucket_trigger" {
  statement_id  = "AllowExecutionFromS3BucketEbayNotifications"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.function.arn
  principal     = "s3.amazonaws.com"
  source_arn    = var.notifications_bucket_arn
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = var.notifications_bucket_id

  lambda_function {
    lambda_function_arn = aws_lambda_function.function.arn
    events              = ["s3:ObjectCreated:*"]
    filter_suffix       = "-ItemListed.xml"
  }

  depends_on = [
    aws_lambda_permission.allow_notification_bucket_trigger
  ]
}