data "aws_vpc" "selected" {
  id = var.vpc_id
}

data "aws_subnets" "transit_gateway_subnets" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }

  filter {
    name   = "tag:TransitGatewayAttachment"
    values = ["true"]
  }
}

data "aws_route_tables" "transit_gateway_route_tables" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }

  filter {
    name   = "tag:TransitGatewayAttachment"
    values = ["true"]
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "tgw_vpc_attachment" {
  subnet_ids         = data.aws_subnets.transit_gateway_subnets.ids
  transit_gateway_id = var.transit_gateway_id
  vpc_id             = var.vpc_id

  tags = {
    Name        = "TransitGatewayAttachment-${data.aws_vpc.selected.tags["Name"]}" # Appends the VPC name
    Environment = "Production"
  }
}

resource "aws_route" "tgw_routes" {
  for_each = toset(data.aws_route_tables.transit_gateway_route_tables.ids)

  route_table_id         = each.value
  destination_cidr_block = "10.0.0.0/8"
  transit_gateway_id     = var.transit_gateway_id
}
