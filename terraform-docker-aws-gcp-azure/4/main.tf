provider "aws" {
  region = "eu-west-3"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "app-vpc"
  cidr = var.vpc_cidr_block

  azs             = [var.availability_zone]
  # private subnets not required, because there is only 1 app-server in this case
  # private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = [var.subnet_cidr_block]

  public_subnet_tags = {
    Name = "${var.env}-subnet-1"
  }
  
  tags = {
    Name = "${var.env}-vpc"
  }
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
  vpc_id               = module.vpc.vpc_id
  subnet_id            = module.vpc.public_subnets[0]
  image_name           = var.image_name
}
