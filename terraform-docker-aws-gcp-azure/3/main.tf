provider "aws" {
  region = "eu-west-3"
}

module "app-subnet" {
  source                 = "./modules/subnet"
  availability_zone      = var.availability_zone
  env                    = var.env
  subnet_cidr_block      = var.subnet_cidr_block
  vpc_id                 = aws_vpc.app-vpc.id
  default_route_table_id = aws_vpc.app-vpc.default_route_table_id
}

module "app-server" {
  source               = "./modules/server"
  availability_zone    = var.availability_zone
  env                  = var.env
  allowed_ips          = var.allowed_ips
  all_ips              = var.all_ips
  instance_type        = var.instance_type
  public_key_location  = var.public_key_location
  private_key_location = var.private_key_location
  vpc_id               = aws_vpc.app-vpc.id
  subnet_id            = module.app-subnet.subnet.id
  image_name           = var.image_name
}

resource "aws_vpc" "app-vpc" {
  cidr_block = var.vpc_cidr_block
  tags = {
    Name : "${var.env}-vpc",
  }
}
