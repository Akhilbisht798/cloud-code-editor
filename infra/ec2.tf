# resource "aws_key_pair" "go_server_key" {
#   key_name = "go-server-key"
#   publi_key = file("~/.ssh/id_rsa.pub")
# }


resource "aws_instance" "go-server" {
  ami                         = var.ami
  instance_type               = var.instance_type
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.public_subnets[0].id
  security_groups             = [aws_security_group.go-server-sg.id]
  iam_instance_profile        = aws_iam_instance_profile.ec2_instance_profile.name
  key_name                    = "loginKeyPair"
  user_data                   = <<-EOF
          #!/bin/bash
          sudo apt update -y
          sudo apt-get install -y git curl

          # Install Go.
          wget https://golang.org/dl/go1.23.1.linux-amd64.tar.gz
          sudo tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
          echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile /home/ubuntu/.bashrc
          source /home/ubuntu/.bashrc
          export PATH=$PATH:/usr/local/go/bin

          echo '${jsonencode(aws_subnet.public_subnets[*].id)}' > /home/ubuntu/subnet_ids.json
          echo '${jsonencode(aws_security_group.socket-server-sg.id)}' > /home/ubuntu/security_group.json

          git clone https://github.com/Akhilbisht798/cloud-code-editor.git /home/ubuntu/code
          cd /home/ubuntu/code/go-server
          sudo -u ubuntu bash -c '
            export PATH=$PATH:/usr/local/go/bin
            export SUBNET_IDS_FILE="/home/ubuntu/subnet_ids.json"
            export SECURITY_GROUP_FILE="/home/ubuntu/security_group.json"
            export APP_ENV="production"
            export BUCKET=""
            export SECRET_KEY=""
            export DB_URL=""
            go run -v ./cmd/main/ > /home/ubuntu/go-server.log 2>&1 &
          '
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
