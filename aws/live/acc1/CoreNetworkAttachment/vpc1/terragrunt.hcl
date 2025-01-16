terraform {
  source = "../../../../terraform-modules/core-network-attachment"
}

inputs = {
  vpc_id             = "vpc-000e5c561a15c50c8"
  transit_gateway_id = "tgw-0af7b4d6adb179480"
  subnet_ids         = ["subnet-02a26f32d4f6828b8", "subnet-00d128f4c64d87399"]
  route_table_ids    = ["rtb-0aebc187547a4f394"]
}