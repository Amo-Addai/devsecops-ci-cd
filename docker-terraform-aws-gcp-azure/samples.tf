# todo: terraform init - setup terraform in current directory 

provider "aws" { // "provider" name
  version = "~> 2.0"
  region  = "us-east-1"
}

resource "aws_vpc" "example" { // "provider_resource" "variable"
  cidr_block = "10.0.0.0/16"   // IPv4 CIDR Block of 
}

provider "kubernetes" {
  config_context_auth_info = "ops"
  config_context_cluster   = "mycluster"
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "sample-namespace"
  }
}

data "aws_vpc" "resource_name" { // query "provider_resource" for "resource_name" data
  default = true
}

resource "aws_subnet" "new_resource_name" {
  vpc_id            = data.aws_vpc.resource_name.id // get .id property of resource_name data from aws_vpc
  cidr_block        = "172.31.48.0/20"              // use sub-set IP from CIDR Block of existing resource_name
  availability_zone = "eu-west-3a"
}

# todo: terraform apply - execute config.tf files; auto-gen's terraform.tfstate (json file with current state)

resource "aws_vpc" "dev-vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name : "development" // key attribute Name; sets resource name in aws
  }
}

resource "aws_vpc" "dev-vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    vpc-env : dev // custom attribute; for arbitrary key-value pairs
  }
}

/** 

# todo: commenting / removing these (with terraform apply) will also remove the resources from aws (if they existed in the previous .tf config)

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

# todo: terrform destroy - removes all config'd resources

# OR: terraform destroy -target provider_resource.resource_name
# eg: terraform destroy -target aws_subnet.dev-subnet-2

# it's best to remove resources, then terraform apply, so that .tf config files will correspond to the current state of cloud resources
# using terraform detroy -target alone will still leave .tf config files inconsistent with the destroyed cloud resources in its current state

*/

# todo: 

/*

terraform plan - preview list to actions to be executed to reach desired cloud state

terraform apply -auto-approve - to auto-respond to t-apply's confirmation question

terraform state - show the current state

*/

output "attribute_name" {
  value = provider_resource.resource_name.attr // attr ~ id
}

/* # TODO: find out how to output multple values this way
output "ids" {
  vpc-id      = aws_vpc.dev-vpc.id
  subnet-1-id = aws_subnet.dev-subnet-1.id
  subnet-2-id = aws_subnet.dev-subnet-2.id
}
*/

