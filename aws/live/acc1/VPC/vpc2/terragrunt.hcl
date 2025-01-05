terraform {
  source = "../../../../terraform-modules/vpc"
}

inputs = {
  vpc_name          = "vpc1"
  cidr_block        = "10.20.0.0/16"
  region            = "us-east-1"
  environment       = "prod"
}