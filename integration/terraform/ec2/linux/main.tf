resource "aws_instance" "integration-test" {
  ami                    = data.aws_ami.latest.id
  instance_type          = var.ec2_instance_type
  key_name               = var.key_name
  iam_instance_profile   = var.iam_instance_profile
  vpc_security_group_ids = var.vpc_security_group_ids
  provisioner "remote-exec" {
    inline = [
      "cloud-init status --wait",
      "echo clone, build, and install agent",
      "git clone ${var.github_repo}",
      "cd amazon-cloudwatch-agent",
      "git reset --hard ${var.github_sha}",
      "aws s3 cp s3://cloudwatch-agent-integration-bucket/integration-test/binary/${var.github_sha}/linux/${var.arc}/${var.binary_name} .",
      "sudo ${var.install_agent}",
      "echo get ssl pem for localstack and export local stack host name",
      "cd ~/amazon-cloudwatch-agent/integration/localstack/ls_tmp",
      "aws s3 cp s3://cloudwatch-agent-integration-bucket/integration-test/ls_tmp/${var.github_sha} . --recursive",
      "cat ${var.ca_cert_path} > original.pem",
      "cat original.pem snakeoil.pem > combine.pem",
      "sudo cp original.pem /opt/aws/amazon-cloudwatch-agent/original.pem",
      "sudo cp combine.pem /opt/aws/amazon-cloudwatch-agent/combine.pem",
      "export LOCAL_STACK_HOST_NAME=${var.local_stack_host_name}",
      "echo run tests with the tag integration, one at a time, and verbose",
      "cd ~/amazon-cloudwatch-agent",
      "make integration-test"
    ]
    connection {
      type        = "ssh"
      user        = var.user
      private_key = var.ssh_key
      host        = self.public_dns
    }
  }
}

data "aws_ami" "latest" {
  most_recent = true
  owners      = ["self"]

  filter {
    name   = "name"
    values = [var.ami]
  }
}