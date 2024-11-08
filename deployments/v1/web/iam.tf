data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "${title(local.function_name)}AssumeLambdaRole"
  description        = "Role for the ebay-listing-ingester to assume"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

data "aws_iam_policy_document" "allow_lambda_logging" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }
}

resource "aws_iam_policy" "function_logging_policy" {
  name        = "Allow${title(local.function_name)}LambdaLoggingPolicy"
  description = "Policy for lambda cloudwatch logging"
  policy      = data.aws_iam_policy_document.allow_lambda_logging.json
}

data "aws_iam_policy_document" "freeagent_auth_ssm_access_policy" {
  statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:GetParametersByPath"
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/freeagent/live/*"
    ]

    effect = "Allow"
  }

  statement {
    actions = [
      "ssm:PutParameter"
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/freeagent/live/access_token",
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/freeagent/live/expires_in",
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/freeagent/live/timestamp"
    ]

    effect = "Allow"
  }

  statement {
    actions = [
      "kms:Decrypt",
      "kms:Encrypt",
      "kms:GenerateDataKey"
    ]

    resources = [
      "*"
    ]

    effect = "Allow"
  }
}

resource "aws_iam_policy" "freeagent_auth_ssm_access_policy" {
  name   = "freeagent-auth-ssm-access-policy"
  policy = data.aws_iam_policy_document.freeagent_auth_ssm_access_policy.json
}

resource "aws_iam_role_policy_attachment" "freeagent_auth_ssm_policy_attachment" {
  role       = aws_iam_role.lambda.id
  policy_arn = aws_iam_policy.freeagent_auth_ssm_access_policy.arn
}

resource "aws_iam_role_policy_attachment" "lambda_logging_policy_attachment" {
  role       = aws_iam_role.lambda.id
  policy_arn = aws_iam_policy.function_logging_policy.arn
}

resource "aws_iam_role_policy_attachment" "function_read_listings_policy_attachment" {
  role       = aws_iam_role.lambda.id
  policy_arn = var.read_listings_policy_arn
}

resource "aws_iam_role_policy_attachment" "function_write_listings_policy_attachment" {
  role       = aws_iam_role.lambda.id
  policy_arn = var.write_listings_policy_arn
}

resource "aws_iam_role_policy_attachment" "function_read_warehouse_policy_attachment" {
  role       = aws_iam_role.lambda.id
  policy_arn = var.allow_read_warehouse_policy_arn
}

resource "aws_iam_role_policy_attachment" "function_write_warehouse_policy_attachment" {
  role       = aws_iam_role.lambda.id
  policy_arn = var.allow_write_warehouse_policy_arn
}

resource "aws_iam_role_policy_attachment" "function_ebay_auth_ssm_access" {
  role       = aws_iam_role.lambda.id
  policy_arn = var.ebay_auth_ssm_access_policy_arn
}
