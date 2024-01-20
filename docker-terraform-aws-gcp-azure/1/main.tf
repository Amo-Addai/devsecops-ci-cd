
provider "aws" {
    region      = "eu-west-3"
    access_key  = "ACCESS_KEY"
    secret_key  = "SECRET"
}

resource "aws_vpc" "dev-vpc" {
    cidr_block  = "10.0.0.0/16"
}

resource "aws_subnet" "dev-subnet-1" {
    vpc_id              = aws_vpc.dev-vpc.id
    cidr_block          = "10.0.10.0/24"
    availability_zone   = "eu-west-3a"
}

data "aws_vpc" "existing_vpc" {
    default = true
}

resource "aws_subnet" "dev-subnet-2" {
    vpc_id              = data.aws_vpc.existing_vpc.id
    cidr_block          = "172.31.48.0/20"
    availability_zone   = "eu-west-3a"
}
