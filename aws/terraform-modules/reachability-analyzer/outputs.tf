output "network_insights_path_id" {
  value = aws_ec2_network_insights_path.source_to_target.id
}

output "network_insights_path_arn" {
  value = aws_ec2_network_insights_path.source_to_target.arn
}

# Outputs for analysis results
output "analysis_status" {
  value = aws_ec2_network_insights_analysis.source_to_target_analysis.status
}

output "analysis_findings" {
  value = aws_ec2_network_insights_analysis.source_to_target_analysis.explanations
}