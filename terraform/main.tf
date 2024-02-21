terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "main" {
  cidr_block = "192.168.0.0/24"
}

resource "aws_subnet" "az_1" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "192.168.0.0/26"
  availability_zone = "us-east-1a"
}

resource "aws_subnet" "az_2" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "192.168.0.64/26"
  availability_zone = "us-east-1b"
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route" "internet" {
  route_table_id         = aws_vpc.main.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.gw.id
}


resource "aws_sns_topic" "emergency" {
  name = "emergency"
}

resource "aws_sns_topic" "signals" {
  name = "signals"
}

resource "aws_sns_topic" "actions" {
  name = "actions"
}


