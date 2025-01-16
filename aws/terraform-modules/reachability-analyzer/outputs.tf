output "network_insights_path_id" {
  value = aws_ec2_network_insights_path.source_to_target.id
}

output "network_insights_path_arn" {
  value = aws_ec2_network_insights_path.source_to_target.arn
}
