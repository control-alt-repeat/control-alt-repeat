resource "aws_s3_bucket" "warehouse" {
  bucket = "control-alt-repeat-warehouse"
}

resource "aws_dynamodb_table" "warehouse" {
  name           = "control-alt-repeat-warehouse"
  billing_mode   = "PAY_PER_REQUEST"

  attribute {
    name = "ID"
    type = "S"
  }

  hash_key = "ID"

  attribute {
    name = "title"
    type = "S"
  }

  attribute {
    name = "shelf"
    type = "S"
  }

  attribute {
    name = "owner"
    type = "S"
  }

  attribute {
    name = "createdAt"
    type = "N"
  }

  attribute {
    name = "updatedAt"
    type = "N"
  }
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
  name        = "AllowReadWarehouseItemsPolicy"
  description = "Policy for lambda reading warehouse"
  policy      = data.aws_iam_policy_document.allow_warehouse_read.json
}

resource "aws_iam_policy" "allow_write_warehouse" {
  name        = "AllowWriteWarehouseItemsPolicy"
  description = "Policy for lambda writing warehouse"
  policy      = data.aws_iam_policy_document.allow_warehouse_write.json
}
