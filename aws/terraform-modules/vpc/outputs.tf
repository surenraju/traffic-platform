output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnets" {
  value = [
    aws_subnet.public_a.id,
    aws_subnet.public_b.id
  ]
}

output "private_subnets" {
  value = [
    aws_subnet.private_a.id,
    aws_subnet.private_b.id
  ]
}

output "storage_subnets" {
  value = [
    aws_subnet.storage_a.id,
    aws_subnet.storage_b.id
  ]
}

output "tgw_attachment_subnets" {
  value = [
    aws_subnet.tgw_attachment_a.id,
    aws_subnet.tgw_attachment_b.id
  ]
}

output "private_route_table_ids" {
   value = [
    aws_route_table.private.id
  ]
}
