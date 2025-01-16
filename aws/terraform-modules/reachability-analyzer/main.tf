terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.83.1"
    }
  }
}

provider "aws" {
  region = "us-east-1" # Replace with your desired region
}

# Provision VM in the source VPC and subnet
resource "aws_instance" "source_vm" {
  ami           = "ami-05576a079321f21f8"
  instance_type = "t2.micro"
  subnet_id     = var.source_subnet_id
  tags = merge(var.tags, { Name = "Source-VM" })
}

# Fetch VPC names for tagging or descriptive purposes
data "aws_vpc" "source_vpc" {
  id = var.source_vpc_id
}

data "aws_vpc" "target_vpc" {
  id = var.target_vpc_id
}

# Fetch random IP from the target subnet
data "aws_subnet" "target_subnet" {
  id = var.target_subnet_id
}

# Use the first IP from the target subnet range (example logic)
locals {
  target_ip = cidrhost(data.aws_subnet.target_subnet.cidr_block, 10)
}

# Network Insights Path: Source VM -> Target Subnet IP
resource "aws_ec2_network_insights_path" "source_to_target" {
  source         = aws_instance.source_vm.arn
  destination_ip = local.target_ip
  destination  = "8.8.8.8"
  protocol       = var.protocol

  tags = merge(
    var.tags,
    {
      Name = format("%s-to-%s-Path", data.aws_vpc.source_vpc.tags["Name"], data.aws_vpc.target_vpc.tags["Name"])
    }
  )
}