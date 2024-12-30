terraform {
  source = "../../terraform-modules/vpc"
}

inputs = {
  vpc_name          = "vpc1"
  cidr_block        = "10.10.0.0/16"
  availability_zones = ["us-east-1a", "us-east-1b"]
  environment       = "prod"
}
