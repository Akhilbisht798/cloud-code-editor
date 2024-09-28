resource "aws_s3_bucket" "project" {
  bucket = "user-project-code-storage-798"
  tags = {
    Name = "Project"
  }
}
