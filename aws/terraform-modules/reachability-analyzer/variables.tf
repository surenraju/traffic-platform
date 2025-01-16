variable "source_vpc_id" {
  description = "The ID of the source VPC"
  type        = string
}

variable "source_subnet_id" {
  description = "The ID of the source subnet"
  type        = string
}

variable "target_vpc_id" {
  description = "The ID of the target VPC"
  type        = string
}

variable "target_subnet_id" {
  description = "The ID of the target subnet"
  type        = string
}

variable "protocol" {
  description = "The protocol for network reachability analysis (e.g., tcp or udp)"
  type        = string
  default     = "tcp"
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}