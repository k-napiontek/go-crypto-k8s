variable "region" {
    type = string
    description = "AWS region to provision infrastructure"
}

variable "bucket_name" {
  type = string

}

variable "github_repos" {
  type = list(string)
}