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
  default = "1024"
}

variable "fargate_memory" {
  type    = string
  default = "2048"
}

variable "network_interface_id" {
  type    = string
  default = "network_id_from_aws"
}

variable "ami" {
  type    = string
  default = "ami-0e86e20dae9224db8"
}

variable "instance_type" {
  type    = string
  default = "t3.medium"
}

variable "cloudwatch_group" {
  type        = string
  default     = "task_logs_cloudwatch"
  description = "cloud watch for ecs tasks"
}
