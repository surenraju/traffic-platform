resource "aws_ec2_transit_gateway_vpc_attachment" "tgw_vpc_attachment" {
  subnet_ids         = var.subnet_ids
  transit_gateway_id = var.transit_gateway_id
  vpc_id             = var.vpc_id

  tags = {
    Name        = "TransitGatewayAttachment-${var.vpc_id}"
    Environment = "Production"
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_route" "tgw_routes" {
  for_each = toset(var.route_table_ids)

  route_table_id         = each.value
  destination_cidr_block = "10.0.0.0/8"
  transit_gateway_id     = var.transit_gateway_id

  lifecycle {
    prevent_destroy = true
  }
}