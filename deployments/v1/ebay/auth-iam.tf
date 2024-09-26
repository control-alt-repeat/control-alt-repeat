data "aws_iam_policy_document" "ebay_auth_ssm_access_policy" {
  statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:GetParametersByPath"
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/ebay/live/*"
    ]

    effect = "Allow"
  }

  statement {
    actions = [ 
      "ssm:PutParameter" 
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/ebay/live/access_token",
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/ebay/live/expires_in"
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

resource "aws_iam_policy" "ebay_auth_ssm_access_policy" {
  name   = "ebay-auth-ssm-access-policy"
  policy = data.aws_iam_policy_document.ebay_auth_ssm_access_policy.json
}
