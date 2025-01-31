provider "aws" {
  region = var.region
}

resource "aws_vpc" "main" {
  cidr_block           = var.cidr_block
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name                     = "${var.vpc_name}"
    Environment              = "prod"
    TransitGatewayAttachment = "true"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "${var.vpc_name}-igw"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "${var.vpc_name}-public-rt"
  }
}

resource "aws_subnet" "public_a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.cidr_block, 8, 1)
  availability_zone       = "${var.region}a"   
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.vpc_name}-public-a"
  }
}

resource "aws_subnet" "public_b" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.cidr_block, 8, 2)
  availability_zone       = "${var.region}b"   
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.vpc_name}-public-b"
  }
}

resource "aws_nat_gateway" "nat_a" {
  allocation_id = aws_eip.nat_a.id
  subnet_id     = aws_subnet.public_a.id

  tags = {
    Name = "${var.vpc_name}-nat-a"
  }
}

resource "aws_nat_gateway" "nat_b" {
  allocation_id = aws_eip.nat_b.id
  subnet_id     = aws_subnet.public_b.id

  tags = {
    Name = "${var.vpc_name}-nat-b"
  }
}

resource "aws_eip" "nat_a" {}

resource "aws_eip" "nat_b" {}

resource "aws_subnet" "private_a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.cidr_block, 8, 3)
  availability_zone = "${var.region}a"
  
  tags = {
    Name = "${var.vpc_name}-private-a"
  }
}

resource "aws_subnet" "private_b" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.cidr_block, 8, 4)
  availability_zone = "${var.region}b"
  
  tags = {
    Name = "${var.vpc_name}-private-b"
  }
}

resource "aws_subnet" "storage_a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.cidr_block, 8, 5)
  availability_zone = "${var.region}a"
  
  tags = {
    Name = "${var.vpc_name}-storage-a"
  }
}

resource "aws_subnet" "storage_b" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.cidr_block, 8, 6)
  availability_zone = "${var.region}b"
  
  tags = {
    Name = "${var.vpc_name}-storage-b"
  }
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat_a.id  # Use nat_a for routing
  }

  tags = {
    Name = "${var.vpc_name}-private-rt"
    TransitGatewayAttachment = "true"
  }
}

resource "aws_route_table_association" "public_a" {
  subnet_id      = aws_subnet.public_a.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "public_b" {
  subnet_id      = aws_subnet.public_b.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private_a" {
  subnet_id      = aws_subnet.private_a.id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "private_b" {
  subnet_id      = aws_subnet.private_b.id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "storage_a" {
  subnet_id      = aws_subnet.storage_a.id
  route_table_id = aws_route_table.private.id
}

resource "aws_route_table_association" "storage_b" {
  subnet_id      = aws_subnet.storage_b.id
  route_table_id = aws_route_table.private.id
}

# Transit Gateway Attachment Subnets
resource "aws_subnet" "tgw_attachment_a" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.cidr_block, 8, 7)  # /28 subnet allocation
  availability_zone = "${var.region}a"
  
  tags = {
    Name                     = "${var.vpc_name}-tgw-attachment-a"
    TransitGatewayAttachment = "true"
  }
}

resource "aws_subnet" "tgw_attachment_b" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.cidr_block, 8, 8)  # /28 subnet allocation
  availability_zone = "${var.region}b"
  
  tags = {
    Name                     = "${var.vpc_name}-tgw-attachment-b"
    TransitGatewayAttachment = "true"
  }
}

# Used for reachability testing
resource "aws_security_group" "reachability_test_sg" {
  name        = "${var.vpc_name}-reachability-test-sg"
  description = "Allow any traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name                  = "${var.vpc_name}-ReachabilityTestSG"
    ReachabilityTestSG    = "true"
  }
}

locals {
  all_reachability_subnets = [
    {
      id   = aws_subnet.private_a.id
      name = "${var.vpc_name}-private-a"
    },
    {
      id   = aws_subnet.private_b.id
      name = "${var.vpc_name}-private-b"
    },
    {
      id   = aws_subnet.storage_a.id
      name = "${var.vpc_name}-storage-a"
    },
    {
      id   = aws_subnet.storage_b.id
      name = "${var.vpc_name}-storage-b"
    }
  ]
}

# Create ENIs for each subnet (Private + Storage),  will be used for reachability testing
resource "aws_network_interface" "reachability_test" {
  for_each        = { for subnet in local.all_reachability_subnets : subnet.name => subnet }
  subnet_id       = each.value.id
  security_groups = [aws_security_group.reachability_test_sg.id]
  depends_on      = [aws_security_group.reachability_test_sg]
  tags = {
    Name                  = "${each.key}-reachability-test-eni"
    ReachabilityTestENI   = "true"
  }
}