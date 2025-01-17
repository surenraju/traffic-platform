variable "source_vpc_id" {
  description = "ID of the source VPC"
  type        = string
}

variable "target_vpc_id" {
  description = "ID of the target VPC"
  type        = string
}

variable "protocol" {
  description = "Protocol for the reachability analyzer (e.g., tcp, udp)"
  type        = string
  default     = "tcp"
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}