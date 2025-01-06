terraform {
  source = "../../../../terraform-modules/core-network-attachment"
}

inputs = {
  vpc_id = "vpc-0fd0f0f20390a839c"
  transit_gateway_id = "tgw-0645ed9c92b0e2a7a"
}
