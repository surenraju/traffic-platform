output "transit_gateway_subnets" {
  description = "List of subnets with TransitGatewayAttachment = 'true'."
  value       = data.aws_subnets.transit_gateway_subnets.ids
}

output "transit_gateway_route_tables" {
  description = "List of route tables with TransitGatewayAttachment = 'true'."
  value       = data.aws_route_tables.transit_gateway_route_tables.ids
}