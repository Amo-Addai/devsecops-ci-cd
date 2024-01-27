
output "aws_ami" {
  value = module.app-server.data-aws_ami.aws_ami.amazon-linux-image-latest
}

output "aws_ami_id" {
  value = module.app-server.data-aws_ami.aws_ami.amazon-linux-image-latest.id
}

output "ec2_public_ip" {
  value = module.app-server.server.public_ip
}
