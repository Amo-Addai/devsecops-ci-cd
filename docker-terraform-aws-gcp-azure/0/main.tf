
provider "aws" {
    region      = "eu-west-3"
    access_key  = "ACCESS_KEY"
    secret_key  = "SECRET"
}

resource "aws_vpc" "dev-vpc" {
    cidr_block  = "10.0.0.0/16"
    tags        = {
        Name: "development"
    }
}

resource "aws_subnet" "dev-subnet-1" {
    vpc_id              = aws_vpc.dev-vpc.id
    cidr_block          = "10.0.10.0/24"
    availability_zone   = "eu-west-3a"
    tags                = {
        Name: "subnet-1-dev"
    }
}