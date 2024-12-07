resource "aws_dynamodb_table" "orders" {
  name           = "control-alt-repeat-orders"
  billing_mode   = "PAY_PER_REQUEST"

  hash_key = "ID"

  attribute {
    name = "ID"
    type = "S"
  }
  
  attribute {
    name = "ItemID"
    type = "S"
  }

  attribute {
    name = "MarketplaceOrderID"
    type = "S"
  }

  attribute {
    name = "MarketplaceSiteID"
    type = "S"
  }

  attribute {
    name = "ShippingCarrier"
    type = "S"
  }

  attribute {
    name = "ShippingTrackingNumber" 
    type = "S"
  }

  attribute {
    name = "ShippingCost"
    type = "S"
  }
}
