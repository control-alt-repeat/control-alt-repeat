resource "aws_s3_bucket" "label_print_buffer" {
  bucket = "control-alt-repeat-label-print-buffer"
}

data "aws_iam_policy_document" "allow_label_print_buffer_read" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.label_print_buffer.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.label_print_buffer.arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "allow_label_print_buffer_write" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.label_print_buffer.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:*Object",
    ]

    resources = [
      "${aws_s3_bucket.label_print_buffer.arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "allow_label_print_buffer_delete_object" {
  statement {
    effect = "Allow"
    actions = [
      "s3:DeleteObject",
    ]

    resources = [
      "${aws_s3_bucket.label_print_buffer.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "allow_read_label_print_buffer" {
  name        = "AllowReadLabelPrintBufferItemsPolicy"
  description = "Policy for lambda reading label_print_buffer"
  policy      = data.aws_iam_policy_document.allow_label_print_buffer_read.json
}

resource "aws_iam_policy" "allow_write_label_print_buffer" {
  name        = "AllowWriteLabelPrintBufferItemsPolicy"
  description = "Policy for lambda writing label_print_buffer"
  policy      = data.aws_iam_policy_document.allow_label_print_buffer_write.json
}

resource "aws_iam_policy" "allow_delete_label_print_buffer" {
  name        = "AllowDeleteLabelPrintBufferItemsPolicy"
  description = "Policy for lambda writing label_print_buffer"
  policy      = data.aws_iam_policy_document.allow_label_print_buffer_delete.json
}

resource "aws_ssm_parameter" "label_printer_host_domain" {
  name  = "/control_alt_repeat/ebay/live/label_printer/host_domain"
  type  = "String"
  value = "not_yet_defined"
}

data "aws_iam_policy_document" "allow_label_print_host_read" {
  statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:GetParametersByPath"
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/ebay/live/label_printer/*"
    ]

    effect = "Allow"
  }

  statement {
    actions = [
      "kms:Decrypt"
    ]

    resources = [
      "*"
    ]

    effect = "Allow"
  }
}

data "aws_iam_policy_document" "allow_label_print_host_write" {
    statement {
    actions = [
      "ssm:GetParameter",
      "ssm:GetParameters",
      "ssm:GetParametersByPath"
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/ebay/live/label_printer/*"
    ]

    effect = "Allow"
  }

  statement {
    actions = [ 
      "ssm:PutParameter" 
    ]

    resources = [
      "arn:aws:ssm:*:*:parameter/control_alt_repeat/ebay/live/label_printer/*"
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

resource "aws_iam_policy" "allow_read_label_print_host" {
  name        = "AllowReadLabelPrintHostItemsPolicy"
  description = "Policy for lambda reading label_print/host"
  policy      = data.aws_iam_policy_document.allow_label_print_host_read.json
}

resource "aws_iam_policy" "allow_write_label_print_host" {
  name        = "AllowWriteLabelPrintHostItemsPolicy"
  description = "Policy for lambda writing label_print/host"
  policy      = data.aws_iam_policy_document.allow_label_print_host_write.json
}

resource "aws_iam_user" "label_printer" {
  name = "label-printer"
  path = "/live/"
}

resource "aws_iam_access_key" "label_printer" {
  user = aws_iam_user.label_printer.name
}

resource "aws_iam_user_policy_attachment" "label_printer_write_host" {
  user       = aws_iam_user.label_printer.name
  policy_arn = aws_iam_policy.allow_write_label_print_host.arn
}

resource "aws_iam_user_policy_attachment" "label_printer_read_print_buffer" {
  user       = aws_iam_user.label_printer.name
  policy_arn = aws_iam_policy.allow_read_label_print_buffer.arn
}

resource "aws_iam_user_policy_attachment" "label_printer_delete_print_buffer" {
  user       = aws_iam_user.label_printer.name
  policy_arn = aws_iam_policy.allow_label_print_buffer_delete_object.arn
}
