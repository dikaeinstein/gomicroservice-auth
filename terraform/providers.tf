terraform {
  backend "s3" {
    bucket = "dikaeinstein-gomicroservice-terraform-state"
    key    = "gomicroservice-auth.tfstate"
    region = "eu-west-2"

    dynamodb_table = "gomicroservice-auth-terraform-state-lock"
    encrypt        = true
  }
}

provider "aws" {
  version = "~> 2.32"
  region  = "eu-west-2"
}

provider "datadog" {}
