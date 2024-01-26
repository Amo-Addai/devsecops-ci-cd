provider "aws" {
  region = "eu-west-3"
}

variable "availability_zone" {}
variable "env" {}
variable "vpc_cidr_block" {}
variable "subnet_cidr_block" {}
variable "allowed_ips" {}
variable "all_ips" {}
variable "instance_type" {}
variable "public_key_location" {}
variable "private_key_location" {}

resource "aws_vpc" "app-vpc" {
  cidr_block = var.vpc_cidr_block
  tags = {
    Name : "${var.env}-vpc",
  }
}

resource "aws_subnet" "app-subnet-1" {
  vpc_id            = aws_vpc.app-vpc.id
  cidr_block        = var.subnet_cidr_block
  availability_zone = var.availability_zone
  tags = {
    Name : "${var.env}-subnet-1"
  }
}

resource "aws_internet_gateway" "app-igw" {
  vpc_id = aws_vpc.app-vpc.id
  tags = {
    Name : "${var.env}-igw"
  }
}

resource "aws_default_route_table" "app-rtb-default" {
  default_route_table_id = aws_vpc.app-vpc.default_route_table_id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.app-igw.id
  }
  tags = {
    Name : "${var.env}-rtb-default"
  }
}

resource "aws_route_table" "app-rtb" {
  vpc_id = aws_vpc.app-vpc.id
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

resource "aws_default_security_group" "app-sg-default" {
  vpc_id = aws_vpc.app-vpc.id
  ingress {
    from_port  = 22
    to_port    = 22
    protocol   = "tcp"
    cidr_block = var.allowed_ips
  }
  ingress {
    from_port  = 8080
    to_port    = 8080
    protocol   = "tcp"
    cidr_block = var.all_ips
  }
  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_block      = var.all_ips
    prefix_list_ids = []
  }
  tags = {
    Name : "${var.env}-sg-default"
  }
}

resource "aws_security_group" "app-sg" {
  name   = "app-sg"
  vpc_id = aws_vpc.app-vpc.id
  ingress {
    from_port  = 22
    to_port    = 22
    protocol   = "tcp"
    cidr_block = var.allowed_ips
  }
  ingress {
    from_port  = 8080
    to_port    = 8080
    protocol   = "tcp"
    cidr_block = var.all_ips
  }
  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_block      = var.all_ips
    prefix_list_ids = []
  }
  tags = {
    Name : "${var.env}-sg"
  }
}

data "aws_ami" "amazon-linux-image-latest" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

output "aws_ami" {
  value = data.aws_ami.amazon-linux-image-latest
}

output "aws_ami_id" {
  value = data.aws_ami.amazon-linux-image-latest.id
}

resource "aws_key_pair" "ssh-key" {
  key_name   = "server-key"
  public_key = file(var.public_key_location)
}

resource "aws_instance" "app-server" {
  ami                         = data.aws_ami.amazon-linux-image-latest.id
  instance_type               = var.instance_type
  subnet_id                   = aws_subnet.app-subnet-1.id
  vpc_security_group_ids      = [aws_default_security_group.app-sg-default.id]
  availability_zone           = var.availability_zone
  associate_public_ip_address = true
  key_name                    = aws_key_pair.ssh-key.key_name

  user_data = file("entrypoint.sh")

  tags = {
    Name = "${var.env}-server"
  }

  provisioner "file" {
    source      = "entrypoint.sh"
    destination = "/path/to/entrypoint.sh"
  }

  provisioner "local-exec" {
    command = "echo ${self.public_ip} > output.txt"
  }

  provisioner "remote-exec" {
    script = file("entrypoint.sh")
  }

  connection {
    type        = "ssh"
    host        = self.public_ip
    user        = "ec2-user"
    private_key = file(var.private_key_location)
  }
}

output "ec2_public_ip" {
  value = aws_instance.app-server.public_ip
}
