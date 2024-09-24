resource "aws_instance" "go-server" {
  ami                         = var.ami
  instance_type               = var.instance_type
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.public_subnets[0].id
  iam_instance_profile        = aws_iam_instance_profile.ec2_instance_profile
  user_data                   = <<-EOF
          #!/bin/bash
          sudo apt update -y
          sudo apt-get install -y git curl

          # Install Go.
          wget https://golang.org/dl/go1.23.1.linux-amd64.tar.gz
          sudo tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
          echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile

          git clone https://github.com/Akhilbisht798/cloud-code-editor.git /home/ubuntu/
          cd /home/ubuntu/cloud-code-editor/go-server

          /usr/local/go/bin/go build -o server .
          nohup ./server &
          EOF
  tags = {
    Name = "Go-Server"
  }
}

resource "aws_security_group" "go-server-sg" {
  name        = "go-server-sg"
  description = "Security group for go server"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 8080 # Adjust this to your Go server's port
    to_port     = 8080
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
    Name = "go-server-sg"
  }
}

output "public_ip" {
  value = aws_instance.go-server.public_ip
}

output "private_ip" {
  value = aws_instance.go-server.private_ip
}
