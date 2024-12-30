terraform {
  source = "../../terraform-modules/vpc"
}

inputs = {
  vpc_name          = "vpc2"
  cidr_block        = "10.20.0.0/16"
  availability_zones = ["us-east-1a", "us-east-1b"]
  environment       = "dev"
}
