resource "aws_s3_bucket" "listings" {
  bucket = "control-alt-repeat-v1-listings"
}

data "aws_iam_policy_document" "allow_listings_read" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.listings.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.listings.arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "allow_listings_write" {
  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.listings.arn
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:*Object",
    ]

    resources = [
      "${aws_s3_bucket.listings.arn}/*",
    ]
  }
}

resource "aws_iam_policy" "allow_read_listings" {
  name        = "AllowReadEbayListingPolicy2"
  description = "Policy for lambda reading listings"
  policy      = data.aws_iam_policy_document.allow_listings_read.json
}

resource "aws_iam_policy" "allow_write_listings" {
  name        = "AllowWriteEbayListingPolicy2"
  description = "Policy for lambda writing listings"
  policy      = data.aws_iam_policy_document.allow_listings_write.json
}
