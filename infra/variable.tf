# variable "aws_access_key" {
#   description = "The IAM public access key"
# }

# variable "aws_secret_key" {
#   description = "IAM secret Key"
# }

variable "aws_region" {
  description = "The AWS region"
  default     = "us-east-1"
}

variable "public_subnets_cidrs" {
  type        = list(string)
  description = "public subnets cidr values"
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "private_subnets_cidrs" {
  type        = list(string)
  description = "private subnets cidr values"
  default     = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]
}

variable "fargate_cpu" {
  type    = string
  default = "256"
}

variable "fargate_memory" {
  type    = string
  default = "512"
}

variable "network_interface_id" {
  type    = string
  default = "network_id_from_aws"
}

variable "ami" {
  type    = string
  default = "ami-06ceb6b6dca8ff42f"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}
