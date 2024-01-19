
provider "aws" {
    version = "~> 2.0"
    region  = "us-east-1"
}

resource "aws_vpc" "example" {
    cidr_block  = "10.0.0.0/16"
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
