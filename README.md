# Terraform Provider for CircleCI Server

A Terraform provider that enables project management for self-hosted CircleCI Server installations using the v1.1 API.

[![Release](https://img.shields.io/github/v/release/anduril/terraform-provider-circleci-server)](https://github.com/anduril/terraform-provider-circleci-server/releases)
[![License](https://img.shields.io/github/license/anduril/terraform-provider-circleci-server)](LICENSE)

## Why This Provider?

**Problem:** CircleCI Server (self-hosted) doesn't support the v2 API's project creation endpoint that the official CircleCI Terraform provider requires. The official provider works with CircleCI Cloud but fails on CircleCI Server installations.

**Solution:** This provider uses the CircleCI v1.1 API's `follow` endpoint, which is supported by CircleCI Server and enables project provisioning through Terraform.

## Features

- ✅ Create CircleCI projects by "following" repositories (v1.1 API)
- ✅ Delete CircleCI projects by "unfollowing" repositories
- ✅ Import existing projects into Terraform state
- ✅ Works with CircleCI Server (self-hosted)
- ✅ Compatible alongside the official CircleCI provider

## Requirements

- Terraform >= 1.0
- Go >= 1.21 (for development)
- CircleCI Server with v1.1 API access
- CircleCI API token with project management permissions

## Installation

### Terraform Registry (Recommended)

```hcl
terraform {
  required_providers {
    circleci-server = {
      source  = "anduril/circleci-server"
      version = "~> 1.0"
    }
  }
}
```

### Local Development

```bash
# Clone and build
git clone https://github.com/anduril/terraform-provider-circleci-server.git
cd terraform-provider-circleci-server
make install
```

## Usage

### Basic Example

```hcl
terraform {
  required_providers {
    circleci-server = {
      source  = "anduril/circleci-server"
      version = "~> 1.0"
    }
  }
}

provider "circleci-server" {
  host  = "https://circleci.example.com"  # Your CircleCI Server URL
  token = var.circleci_token               # CircleCI API token
}

# Follow a project to enable CircleCI builds
resource "circleci-server_project_follow" "my_app" {
  vcs_type     = "github"      # or "bitbucket"
  organization = "my-org"
  project      = "my-repo"
}
```

### Using with Official CircleCI Provider

This provider handles project creation; use the official provider for everything else:

```hcl
terraform {
  required_providers {
    circleci-server = {
      source  = "anduril/circleci-server"
      version = "~> 1.0"
    }
    circleci = {
      source  = "CircleCI-Public/circleci"
      version = "~> 0.3"
    }
  }
}

provider "circleci-server" {
  host  = "https://circleci.example.com"
  token = var.circleci_token
}

provider "circleci" {
  host = "https://circleci.example.com"
  key  = var.circleci_token
}

# Create project with this provider
resource "circleci-server_project_follow" "app" {
  vcs_type     = "github"
  organization = "my-org"
  project      = "my-app"
}

# Manage project settings with official provider
resource "circleci_environment_variable" "api_key" {
  project = "gh/my-org/my-app"
  name    = "API_KEY"
  value   = var.api_key

  depends_on = [circleci-server_project_follow.app]
}
```

### Importing Existing Projects

```bash
terraform import circleci-server_project_follow.app github/my-org/my-app
```

## Resources

### `circleci-server_project_follow`

Follows a CircleCI project, enabling builds for the repository.

**Arguments:**

- `vcs_type` (Required, String) - VCS type: `"github"` or `"bitbucket"`
- `organization` (Required, String) - Organization or username
- `project` (Required, String) - Repository/project name

**Attributes:**

- `id` (String) - Resource ID in format `{vcs_type}/{organization}/{project}`

**Import:**

```bash
terraform import circleci-server_project_follow.example github/org-name/repo-name
```

## Development

```bash
# Install dependencies
go mod tidy

# Build
make build

# Install locally
make install

# Run tests
make test

# Format code
make fmt
```

## Limitations

- Only supports project creation/deletion (follow/unfollow operations)
- Read operations assume the project exists after creation
- Designed specifically for CircleCI Server; may not work with CircleCI Cloud
- Requires CircleCI Server v1.1 API support

## Contributing

Contributions are welcome! Please open an issue or pull request.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Credits

Built with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).
