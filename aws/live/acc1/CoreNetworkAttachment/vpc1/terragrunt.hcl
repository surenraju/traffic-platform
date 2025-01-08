terraform {
  source = "../../../../terraform-modules/core-network-attachment"
}

inputs = {
  vpc_id             = "vpc-0cdf83acc85ee7d40"
  transit_gateway_id = "tgw-0293d559f0df7aaae"
  subnet_ids         = ["subnet-08c8921c55256136e", "subnet-0b7063b0b4cecdaa5"]
  route_table_ids    = ["rtb-0ba7aed0ee05a4a77"]
}