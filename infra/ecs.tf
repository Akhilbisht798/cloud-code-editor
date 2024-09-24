resource "aws_ecs_cluster" "worker" {
  name = "socket-servers"
}

resource "aws_ecs_task_definition" "socket-task-defination" {
  family                   = "socket-server-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.fargate_cpu
  memory                   = var.fargate_memory
  container_definitions = jsonencode([
    {
      name      = "socket-server"
      image     = "akhilbisht798/socket-server"
      essential = true
      environment = [
        {
          name  = "SERVER_URL"
          value = aws_instance.go-server.private_ip
        }
      ]
      portMappings = [
        {
          containerPort = 5000
          hostPort      = 5000
          protocol      = "tcp"
        },
        {
          containerPort = 3000
          hostPort      = 3000
          protocol      = "tcp"
        }
      ]
    }
  ])
}

resource "aws_security_group" "socket-server-sg" {
  name        = "socket-server-sg"
  description = "sg for socket server"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 5000
    to_port     = 5000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "socket-server-sg"
  }
}
