resource "aws_s3_bucket" "ebay_incoming_notifications" {
  bucket = "control-alt-repeat-live-ebay-incoming-notifications"
}

resource "aws_s3_bucket_lifecycle_configuration" "ebay_incoming_notifications_lifecycle_1" {
  bucket = aws_s3_bucket.ebay_incoming_notifications.id
  rule {
    status = "Enabled"
    id     = "expire_all_files"
    expiration {
        days = 7
    }
  }
}

data "aws_iam_policy_document" "allow_ebay_incoming_notifications_read" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.ebay_incoming_notifications.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.ebay_incoming_notifications.arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "allow_ebay_incoming_notifications_write" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.ebay_incoming_notifications.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:*Object",
    ]

    resources = [
      "${aws_s3_bucket.ebay_incoming_notifications.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "allow_read_ebay_incoming_notifications" {
  name        = "AllowReadLabelPrintBufferItemsPolicy"
  description = "Policy for lambda reading ebay_incoming_notifications"
  policy      = data.aws_iam_policy_document.allow_ebay_incoming_notifications_read.json
}

resource "aws_iam_policy" "allow_write_ebay_incoming_notifications" {
  name        = "AllowWriteLabelPrintBufferItemsPolicy"
  description = "Policy for lambda writing ebay_incoming_notifications"
  policy      = data.aws_iam_policy_document.allow_ebay_incoming_notifications_write.json
}
