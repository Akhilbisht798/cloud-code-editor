resource "aws_s3_bucket" "project" {
  bucket = "project"
  tags = {
    Name = "Project"
  }
}
