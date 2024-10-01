resource "aws_s3_bucket" "warehouse" {
  bucket = "control-alt-repeat-warehouse"
}

data "aws_iam_policy_document" "allow_warehouse_read" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.warehouse.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.warehouse.arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "allow_warehouse_write" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.warehouse.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:*Object",
    ]

    resources = [
      "${aws_s3_bucket.warehouse.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "allow_read_warehouse" {
  name        = "AllowReadEbayListingPolicy"
  description = "Policy for lambda reading warehouse"
  policy      = data.aws_iam_policy_document.allow_warehouse_read.json
}

resource "aws_iam_policy" "allow_write_warehouse" {
  name        = "AllowWriteEbayListingPolicy"
  description = "Policy for lambda writing warehouse"
  policy      = data.aws_iam_policy_document.allow_warehouse_write.json
}
