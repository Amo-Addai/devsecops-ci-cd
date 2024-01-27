resource "aws_default_security_group" "app-sg-default" {
  vpc_id = var.vpc_id
  
  ingress {
    from_port  = 22
    to_port    = 22
    protocol   = "tcp"
    cidr_block = var.allowed_ips
  }
  ingress {
    from_port  = 8080
    to_port    = 8080
    protocol   = "tcp"
    cidr_block = var.all_ips
  }
  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_block      = var.all_ips
    prefix_list_ids = []
  }
  tags = {
    Name : "${var.env}-sg-default"
  }
}

resource "aws_security_group" "app-sg" {
  name   = "app-sg"
  vpc_id = var.vpc_id

  ingress {
    from_port  = 22
    to_port    = 22
    protocol   = "tcp"
    cidr_block = var.allowed_ips
  }
  ingress {
    from_port  = 8080
    to_port    = 8080
    protocol   = "tcp"
    cidr_block = var.all_ips
  }
  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_block      = var.all_ips
    prefix_list_ids = []
  }
  tags = {
    Name : "${var.env}-sg"
  }
}

data "aws_ami" "amazon-linux-image-latest" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "name"
    values = [var.image_name]
  }
  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_key_pair" "ssh-key" {
  key_name   = "server-key"
  public_key = file(var.public_key_location)
}

resource "aws_instance" "app-server" {
  ami                         = data.aws_ami.amazon-linux-image-latest.id
  instance_type               = var.instance_type
  subnet_id                   = var.subnet_id
  vpc_security_group_ids      = [aws_default_security_group.app-sg-default.id, aws_security_group.app-sg.id]
  availability_zone           = var.availability_zone
  associate_public_ip_address = true
  key_name                    = aws_key_pair.ssh-key.key_name

  user_data = file("../../entrypoint.sh")

  tags = {
    Name = "${var.env}-server"
  }

  provisioner "file" {
    source      = "../../entrypoint.sh"
    destination = "/path/to/entrypoint.sh"
  }

  provisioner "local-exec" {
    command = "echo ${self.public_ip} > output.txt"
  }

  provisioner "remote-exec" {
    script = file("../../entrypoint.sh")
  }

  connection {
    type        = "ssh"
    host        = self.public_ip
    user        = "ec2-user"
    private_key = file(var.private_key_location)
  }
}
