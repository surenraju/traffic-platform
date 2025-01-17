terraform {
  source = "../../../../terraform-modules/core-network-attachment"
}

inputs = {
  vpc_id             = "vpc-084d0d2e74a789756"
  transit_gateway_id = "tgw-008a46bb27ebdc169"
  subnet_ids         = ["subnet-0c7d54d3c3fb8564a", "subnet-033c5dc62f3f97d10"]
  route_table_ids    = ["rtb-043b544a5533bc2df"]
}