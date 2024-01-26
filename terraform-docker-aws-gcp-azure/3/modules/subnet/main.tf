resource "aws_subnet" "app-subnet-1" {
  vpc_id            = var.vpc_id
  cidr_block        = var.subnet_cidr_block
  availability_zone = var.availability_zone
  tags = {
    Name : "${var.env}-subnet-1"
  }
}

resource "aws_internet_gateway" "app-igw" {
  vpc_id = var.vpc_id
  tags = {
    Name : "${var.env}-igw"
  }
}

resource "aws_default_route_table" "app-rtb-default" {
  default_route_table_id = var.default_route_table_id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.app-igw.id
  }
  tags = {
    Name : "${var.env}-rtb-default"
  }
}

resource "aws_route_table" "app-rtb" {
  vpc_id = var.vpc_id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.app-igw.id
  }
  tags = {
    Name : "${var.env}-rtb"
  }
}

resource "aws_route_table_association" "app-rtb-subnet" {
  subnet_id      = aws_subnet.app-subnet-1.id
  route_table_id = aws_route_table.app-rtb.id
}