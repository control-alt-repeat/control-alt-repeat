resource "aws_s3_bucket" "tokens" {
  bucket = "control-alt-repeat-v1-oauth2-tokens"
}

data "aws_iam_policy_document" "allow_tokens_read" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.tokens.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.tokens.arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "allow_tokens_write" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.tokens.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:*Object",
    ]

    resources = [
      "${aws_s3_bucket.tokens.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "function_read_tokens" {
  name        = "AllowReadEbayAuthTokensPolicy"
  description = "Policy for lambda reading tokens"
  policy      = data.aws_iam_policy_document.allow_tokens_read.json
}

resource "aws_iam_policy" "function_write_tokens" {
  name        = "AllowWriteEbayAuthTokensPolicy"
  description = "Policy for lambda writing tokens"
  policy      = data.aws_iam_policy_document.allow_tokens_write.json
}
