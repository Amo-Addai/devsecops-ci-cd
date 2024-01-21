provider "aws" {
  region     = "eu-west-3"
}

variable availability_zone {} // "eu-west-3a"

variable "environment" {
  description = "current environment"
}

variable "cidr_blocks" {
  description = "vpc cidr blocks"
  type        = list(object({
    cidr_block = string
    name       = string
  }))
}

/*
variable "vpc_cidr_block" {
  description = "vpc cidr block"
  default     = "10.0.10.0/24"
  type        = string
}

variable "subnet_cidr_block" {
  description = "subnet cidr block"
}
*/

resource "aws_vpc" "dev-vpc" {
  cidr_block = var.cidr_blocks[0].cidr_block
  tags = {
    Name : var.cidr_blocks[0].name,
    vpc_env : var.environment
  }
}

resource "aws_subnet" "dev-subnet-1" {
  vpc_id            = aws_vpc.dev-vpc.id
  cidr_block        = var.cidr_blocks[1].cidr_block
  availability_zone = var.availability_zone
  tags = {
    Name : var.cidr_blocks[1].name
  }
}

