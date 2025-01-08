terraform {
  source = "../../../../terraform-modules/core-network-attachment"
}

inputs = {
  vpc_id             = "vpc-05258717303d7775f"
  transit_gateway_id = "tgw-0293d559f0df7aaae"
  subnet_ids         = ["subnet-0bc45b0b0692d58f1", "subnet-053d742645095da82"]
  route_table_ids    = ["rtb-06f02cafd35b25d79"]
}