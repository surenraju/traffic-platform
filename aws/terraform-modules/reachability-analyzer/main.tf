# Fetch source VPC information
data "aws_vpc" "source_vpc" {
  id = var.source_vpc_id
}

# Fetch target VPC information
data "aws_vpc" "target_vpc" {
  id = var.target_vpc_id
}

# Fetch all ENIs in the source VPC with the specified tag
data "aws_network_interfaces" "source_enis" {
  filter {
    name   = "vpc-id"
    values = [var.source_vpc_id]
  }

  filter {
    name   = "tag:ReachabilityTestENI"
    values = ["true"]
  }
}

# Fetch all ENIs in the target VPC with the specified tag
data "aws_network_interfaces" "target_enis" {
  filter {
    name   = "vpc-id"
    values = [var.target_vpc_id]
  }

  filter {
    name   = "tag:ReachabilityTestENI"
    values = ["true"]
  }
}

# Validate that ENIs exist and get their ARNs
locals {
  source_eni_id = length(data.aws_network_interfaces.source_enis.ids) > 0 ? data.aws_network_interfaces.source_enis.ids[0] : null
  target_eni_id = length(data.aws_network_interfaces.target_enis.ids) > 0 ? data.aws_network_interfaces.target_enis.ids[0] : null
}

# Fetch details of the source ENI
data "aws_network_interface" "source_eni" {
  id = local.source_eni_id
}

# Fetch details of the target ENI
data "aws_network_interface" "target_eni" {
  id = local.target_eni_id
}

# Create the Network Insights Path
resource "aws_ec2_network_insights_path" "source_to_target" {
  source      = data.aws_network_interface.source_eni.arn
  destination = data.aws_network_interface.target_eni.arn
  protocol    = var.protocol

  tags = merge(
    var.tags,
    {
      Name = format("%s-to-%s-Path", data.aws_vpc.source_vpc.tags["Name"], data.aws_vpc.target_vpc.tags["Name"])
    }
  )
}

# Null resource to force trigger on every apply
resource "null_resource" "always_run" {
  triggers = {
    timestamp = timestamp()
  }
}

# Run Network Insights Analysis on every Terraform apply
resource "aws_ec2_network_insights_analysis" "source_to_target_analysis" {
  network_insights_path_id = aws_ec2_network_insights_path.source_to_target.id

  tags = merge(
    var.tags,
    {
      Name = format("Analysis-%s-to-%s", data.aws_vpc.source_vpc.tags["Name"], data.aws_vpc.target_vpc.tags["Name"])
    }
  )

  # Always trigger the analysis
  lifecycle {
    replace_triggered_by = [
      null_resource.always_run
    ]
  }
}