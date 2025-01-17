terraform {
  source = "../../../../terraform-modules/core-network-attachment"
}

inputs = {
  vpc_id             = "vpc-0ef5790be12fcd7db"
  transit_gateway_id = "tgw-008a46bb27ebdc169"
  subnet_ids         = ["subnet-074e960b17fe51ef9", "subnet-08b8a9fac4d86ae39"]
  route_table_ids    = ["rtb-0cea0e560dbf9cd2f"]
}