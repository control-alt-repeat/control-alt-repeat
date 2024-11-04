resource "aws_s3_bucket" "warehouse" {
  bucket = "control-alt-repeat-warehouse"
}

resource "aws_dynamodb_table" "warehouse" {
  name           = "control-alt-repeat-warehouse"
  billing_mode   = "PAY_PER_REQUEST"

  hash_key = "ID"

  attribute {
    name = "ID"
    type = "S"
  }

  attribute {
    name = "Shelf"
    type = "S"
  }

  attribute {
    name = "Owner"
    type = "S"
  }

  attribute {
    name = "CreatedAt"
    type = "N"
  }

  attribute {
    name = "UpdatedAt"
    type = "N"
  }

  global_secondary_index {
    name            = "GSI_Shelf"
    hash_key        = "Shelf"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "GSI_Owner"
    hash_key        = "Owner"
    projection_type = "ALL"
  }

  global_secondary_index {
    name               = "CreatedAtIndex"
    hash_key           = "ID"
    range_key          = "CreatedAt"
    projection_type    = "ALL"
  }

  global_secondary_index {
    name               = "UpdatedAtIndex"
    hash_key           = "ID"
    range_key          = "UpdatedAt"
    projection_type    = "ALL"
  }
}

data "aws_iam_policy_document" "allow_warehouse_read" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.warehouse.arn,
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

    statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:BatchGetItem",
      "dynamodb:Query",
      "dynamodb:Scan"
    ]

    resources = [
      aws_dynamodb_table.warehouse.arn,
      "${aws_dynamodb_table.warehouse.arn}/index/*",
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
      aws_s3_bucket.warehouse.arn,
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem",
      "dynamodb:BatchGetItem",
      "dynamodb:BatchWriteItem",
      "dynamodb:Query",
      "dynamodb:Scan"
    ]

    resources = [
      aws_dynamodb_table.warehouse.arn,
      "${aws_dynamodb_table.warehouse.arn}/index/*",
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
