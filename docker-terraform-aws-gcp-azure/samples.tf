# TODO: terraform init - setup terraform in current directory 

provider "aws" { // "provider" name
    version = "~> 2.0"
    region  = "us-east-1"
}

resource "aws_vpc" "example" { // "provider_resource" "variable"
    cidr_block  = "10.0.0.0/16" // IPv4 CIDR Block of 
}

provider "kubernetes" {
    config_context_auth_info    = "ops"
    config_context_cluster      = "mycluster"
}

resource "kubernetes_namespace" "example" {
    metadata {
        name    = "sample-namespace"
    }
}

data "aws_vpc" "resource_name" { // query "provider_resource" for "resource_name" data
    default = true
}

resource "aws_subnet" "new_resource_name" {
    vpc_id              = data.aws_vpc.resource_name.id // get .id property of resource_name data from aws_vpc
    cidr_block          = "172.31.48.0/20" // use sub-set IP from CIDR Block of existing resource_name
    availability_zone   = "eu-west-3a"
}

# TODO: terraform apply - execute config.tf files
