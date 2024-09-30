# resource "aws_key_pair" "go_server_key" {
#   key_name = "go-server-key"
#   publi_key = file("~/.ssh/id_rsa.pub")
# }


resource "aws_instance" "go-server" {
  ami                         = var.ami
  instance_type               = var.instance_type
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.public_subnets[0].id
  iam_instance_profile        = aws_iam_instance_profile.ec2_instance_profile.name
  key_name                    = "loginKeyPair"
  user_data                   = <<-EOF
          #!/bin/bash
          sudo apt update -y
          sudo apt-get install -y git curl

          # Install Go.
          wget https://golang.org/dl/go1.23.1.linux-amd64.tar.gz
          sudo tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
          export PATH=$PATH:/usr/local/go/bin

          echo '{jsonencode(aws_subnet.public_subnets[*].id)}' > /home/ubuntu/subnet_ids.json

          git clone https://github.com/Akhilbisht798/cloud-code-editor.git /home/ubuntu/code
          cd /home/ubuntu/code/go-server
          cat <<EOT > /home/ubuntu/run_server.sh

          #!/bin/bash
          export SUBNET_IDS_FILE="/home/ubuntu/subnet_ids.json"
          export APP_ENV="production"
          export BUCKET="user-project-code-storage-798"
          export SECRET_KEY="secret"

          cd /home/ubuntu/code/go-server
          go run ./cmd/main/
          EOT
          chomd +x /home/ubuntu/code/run_server.sh
          sudo -u ubuntu /home/ubuntu/run_server.sh &
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
    from_port   = 22 # Adjust this to your Go server's port
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

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
