# resource "null_resource" "function_binary" {
#   provisioner "local-exec" {
#     command = "GOOS=linux GOARC=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -tags lambda.norpc -mod=readonly -o ${local.binary_path} ${local.src_path}"
#   }
# }

# data "archive_file" "function_archive" {
#   depends_on = [null_resource.function_binary]

#   type        = "zip"
#   source_file = local.binary_path
#   output_path = local.archive_path
# }

resource "aws_lambda_function" "function" {
  function_name    = "control-alt-repeat-poc"
  description      = "Indexes eBay listings into Control Alt Repeats asset inventory S3 bucket"
  role             = aws_iam_role.lambda.arn
  handler          = "bootstrap"
  filename         = local.archive_path
  source_code_hash = filebase64sha256(local.archive_path)
  runtime          = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.function.function_name}"
  retention_in_days = 7
}
