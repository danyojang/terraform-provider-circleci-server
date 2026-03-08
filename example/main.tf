terraform {
  required_providers {
    circleci-server = {
      source  = "anduril/circleci-server"
      version = "1.0.0"
    }
  }
}

provider "circleci-server" {
  host  = "https://cci.anduril.dev"
  token = var.circleci_token
}

variable "circleci_token" {
  type      = string
  sensitive = true
}

# Example: Follow a project
resource "circleci-server_project_follow" "example" {
  vcs_type     = "github"
  organization = "test-org"
  project      = "my-repo-name"
}

output "project_id" {
  value = circleci-server_project_follow.example.id
}
