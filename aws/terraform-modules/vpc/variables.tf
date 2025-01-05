variable "vpc_name" {
  description = "Name of the VPC"
  type        = string
}

variable "cidr_block" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "region" {
  description = "The AWS region for the VPC"
  type        = string
  default     = "us-east-1" 
}


variable "environment" {
  description = "Environment for the VPC (e.g., dev, prod)"
  type        = string
}