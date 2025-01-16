terraform {
  source = "../../../../terraform-modules/core-network-attachment"
}

inputs = {
  vpc_id             = "vpc-03e1c50b58a38b839"
  transit_gateway_id = "tgw-0af7b4d6adb179480"
  subnet_ids         = ["subnet-09c54f61a9907f180", "subnet-0dbe65a1354c9a757"]
  route_table_ids    = ["rtb-08c4032574fd13d20"]
}