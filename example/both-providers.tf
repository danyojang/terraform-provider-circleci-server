# Example showing how to use both providers together

terraform {
  required_providers {
    # Custom provider for project creation via v1.1 follow API
    circleci-server = {
      source  = "anduril/circleci-server"
      version = "1.0.0"
    }
    # Official provider for everything else (env vars, schedules, etc.)
    circleci = {
      source  = "CircleCI-Public/circleci"
      version = "~> 0.3"
    }
  }
}

provider "circleci-server" {
  host  = "https://cci.anduril.dev"
  token = var.circleci_token
}

provider "circleci" {
  host = "https://cci.anduril.dev"
  key  = var.circleci_token
}

variable "circleci_token" {
  type      = string
  sensitive = true
}

# Step 1: Follow/create the project using custom provider
resource "circleci-server_project_follow" "my_project" {
  vcs_type     = "github"
  organization = "test-org"
  project      = "my-repo"
}

# Step 2: Use official provider to manage project settings
# (These resources will work once v2 API is available for other operations)

# Example: Add environment variable (if official provider supports it)
# resource "circleci_environment_variable" "api_key" {
#   project = "gh/test-org/my-repo"
#   name    = "API_KEY"
#   value   = "secret-value"
#
#   depends_on = [circleci-server_project_follow.my_project]
# }

# Example: Add checkout key (if official provider supports it)
# resource "circleci_checkout_key" "deploy_key" {
#   project = "gh/test-org/my-repo"
#   type    = "deploy-key"
#
#   depends_on = [circleci-server_project_follow.my_project]
# }
