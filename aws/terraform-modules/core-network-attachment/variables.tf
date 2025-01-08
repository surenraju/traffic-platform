variable "vpc_id" {
  description = "The ID of the AWS VPC."
  type        = string
}

variable "transit_gateway_id" {
  description = "The ID of the Transit Gateway."
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for the Transit Gateway attachment."
  type        = list(string)
}

variable "route_table_ids" {
  description = "List of route table IDs for Transit Gateway routing."
  type        = list(string)
}