variable "description" {
  description = "The description of the transit gateway"
  type        = string
  default     = "Transit Gateway to enable cross VPC communication"
}

variable "amazon_side_asn" {
  description = "The ASN for the Amazon side of the transit gateway"
  type        = number
  default     = 64512
}

variable "name" {
  description = "The name of the transit gateway"
  type        = string
}

