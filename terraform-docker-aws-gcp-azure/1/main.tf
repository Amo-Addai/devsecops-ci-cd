provider "aws" {
  region     = "eu-west-3"
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

variable aws_access_key {}
variable aws_secret_key {}
variable availability_zone {}

variable "environment" {
  description = "current environment"
}

variable "vpc_cidr_block" {
  description = "vpc cidr block"
  default     = "10.0.10.0/24"
  type        = string
}

variable "subnet_cidr_block" {
  description = "subnet cidr block"
}

resource "aws_vpc" "dev-vpc" {
  cidr_block = var.vpc_cidr_block
  tags = {
    Name : var.environment,
    vpc_env : "dev"
  }
}

resource "aws_subnet" "dev-subnet-1" {
  vpc_id            = aws_vpc.dev-vpc.id
  cidr_block        = var.subnet_cidr_block
  availability_zone = "eu-west-3a"
  tags = {
    Name : "subnet-1-dev"
  }
}

data "aws_vpc" "existing_vpc" {
  default = true
}

resource "aws_subnet" "dev-subnet-2" {
  vpc_id            = data.aws_vpc.existing_vpc.id
  cidr_block        = "172.31.42.0/20"
  availability_zone = "eu-west-3a"
  tags = {
    Name : "subnet-2-default"
  }
}

output "dev-vpc-id" {
  value = aws_vpc.dev-vpc.id
}

output "dev-subnet-1-id" {
  value = aws_subnet.dev-subnet-1.id
}

output "dev-subnet-2-id" {
  value = aws_subnet.dev-subnet-2.id
}

