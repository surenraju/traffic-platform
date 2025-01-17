terraform {
  source = "../../../../terraform-modules/reachability-analyzer"
}


inputs = {
  source_vpc_id     = "vpc-0ef5790be12fcd7db"
  target_vpc_id     = "vpc-084d0d2e74a789756"
  protocol          = "tcp"
}