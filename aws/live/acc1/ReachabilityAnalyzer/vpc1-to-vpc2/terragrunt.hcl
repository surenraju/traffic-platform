terraform {
  source = "../../../../terraform-modules/reachability-analyzer"
}

dependency "vpc1" {
  config_path = "../../VPC/vpc1"
}

dependency "vpc2" {
  config_path = "../../VPC/vpc2"
}

inputs = {
  source_vpc_id     = dependency.vpc1.outputs.vpc_id
  source_subnet_id  = dependency.vpc1.outputs.private_subnets[0] # First private subnet in vpc1
  target_vpc_id     = dependency.vpc2.outputs.vpc_id
  target_subnet_id  = dependency.vpc2.outputs.private_subnets[0] # First private subnet in vpc2
  protocol          = "tcp"
  tags = {
    Environment = "prod"
    Project     = "Reachability-vpc1-to-vpc2"
  }
}