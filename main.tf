provider "aws" {
  region                  = "eu-west-2"  # You can set this to any AWS region
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  access_key = "abc"
  secret_key = "def"
  endpoints {
    dynamodb = "http://localhost:4566" # Use the LocalStack endpoint
  }
}

resource "aws_dynamodb_table" "this" {
  name           = "cicd-audit"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "Component"
  range_key = "Environment"

  attribute {
    name = "Component"
    type = "S"
  }

  attribute {
    name = "Environment"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "S"
  }

  attribute {
    name = "State"
    type = "S"
  }

  global_secondary_index {
    name = "GSI1"
    hash_key = "State"
    range_key = "Timestamp"
    projection_type = "INCLUDE"
    non_key_attributes = ["Component", "Environment"]
  }
}