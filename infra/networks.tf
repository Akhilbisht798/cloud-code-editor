data "aws_availability_zones" "available" {
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id
  tags = {
    Name = "Project VPC IG"
  }
}

resource "aws_subnet" "public_subnets" {
  count             = length(var.public_subnets_cidrs)
  vpc_id            = aws_vpc.main.id
  cidr_block        = element(var.public_subnets_cidrs, count.index)
  availability_zone = element(data.aws_availability_zones.available, count.index)

  tags = {
    Name = "Public Subnet ${count.index + 1}"
  }
}

resource "aws_subnet" "private_subnets" {
  count             = length(var.private_subnets_cidrs)
  vpc_id            = aws_vpc.main.id
  cidr_block        = element(var.private_subnets_cidrs, count.index)
  availability_zone = element(data.aws_availability_zones.available, count.index)

  tags = {
    Name = "Private Subnet ${count.index + 1}"
  }
}

# routing.
resource "aws_route_table" "second_rt" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "2nd route table"
  }
}

resource "aws_route_table_association" "public_subnet_asso" {
  count          = length(var.public_subnets_cidrs)
  subnet_id      = element(aws_subnet.public_subnets[*].id, count.index)
  route_table_id = aws_route_table.second_rt.id
}


