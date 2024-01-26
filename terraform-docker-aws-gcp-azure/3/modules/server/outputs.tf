
output "server" {
    value = aws_instance.app-server
}

output "data-aws_ami" {
    value = data.aws_ami
}
