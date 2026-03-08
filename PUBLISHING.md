# Publishing to Terraform Registry

This guide walks through publishing the provider to the public Terraform Registry.

## Prerequisites

1. **GitHub Account** - Needs to be a member of the Anduril GitHub organization
2. **GPG Key** - For signing releases
3. **Terraform Registry Account** - Sign up at https://registry.terraform.io

## Step 1: Create Public GitHub Repository

```bash
cd terraform-provider-circleci-server

# Initialize git if not already done
git init

# Create repository on GitHub at: github.com/anduril/terraform-provider-circleci-server
# Then add remote:
git remote add origin https://github.com/anduril/terraform-provider-circleci-server.git

# Add all files
git add .
git commit -m "Initial provider release"
git push -u origin main
```

## Step 2: Generate GPG Key

```bash
# Generate GPG key (if you don't have one)
gpg --full-generate-key
# Choose: RSA and RSA, 4096 bits, no expiration
# Use your work email

# List keys to get the key ID
gpg --list-secret-keys --keyid-format=long
# Output shows: sec   rsa4096/YOUR_KEY_ID 2026-03-08

# Export public key for Terraform Registry
gpg --armor --export YOUR_KEY_ID > gpg-public-key.asc

# Export private key for GitHub Secrets
gpg --armor --export-secret-keys YOUR_KEY_ID > gpg-private-key.asc
```

## Step 3: Add GPG Key to Terraform Registry

1. Go to https://registry.terraform.io
2. Sign in with GitHub
3. Go to Settings → Signing Keys
4. Click "Add GPG Public Key"
5. Paste contents of `gpg-public-key.asc`
6. Save

## Step 4: Add Secrets to GitHub Repository

1. Go to GitHub repository settings
2. Navigate to: Settings → Secrets and variables → Actions
3. Add two secrets:
   - `GPG_PRIVATE_KEY`: Contents of `gpg-private-key.asc`
   - `PASSPHRASE`: Your GPG key passphrase

## Step 5: Publish Provider to Terraform Registry

1. Go to https://registry.terraform.io
2. Click "Publish" → "Provider"
3. Select your GitHub repository: `anduril/terraform-provider-circleci-server`
4. The registry will verify:
   - Repository name matches pattern `terraform-provider-*`
   - Has a valid LICENSE file
   - Has releases with proper assets

## Step 6: Create Your First Release

```bash
# Tag and push a release
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will automatically:
# 1. Build binaries for all platforms
# 2. Sign the release
# 3. Create GitHub release
# 4. Terraform Registry will detect it

# Check progress:
# - GitHub: Actions tab
# - Terraform Registry: Should show v1.0.0 after ~5 minutes
```

## Step 7: Verify Publication

After release completes:

1. Check Terraform Registry: https://registry.terraform.io/providers/anduril/circleci-server
2. Test installation:
   ```bash
   mkdir test-provider
   cd test-provider
   cat > main.tf <<EOF
   terraform {
     required_providers {
       circleci-server = {
         source  = "anduril/circleci-server"
         version = "~> 1.0"
       }
     }
   }
   EOF
   terraform init
   ```

## Troubleshooting

### Release fails with GPG error
- Verify `GPG_PRIVATE_KEY` and `PASSPHRASE` secrets are correct
- Test locally: `gpg --armor --export-secret-keys YOUR_KEY_ID`

### Terraform Registry doesn't detect release
- Ensure tag starts with `v` (e.g., `v1.0.0`)
- Check release has these assets:
  - Multiple `.zip` files (one per platform)
  - `terraform-provider-circleci-server_1.0.0_SHA256SUMS`
  - `terraform-provider-circleci-server_1.0.0_SHA256SUMS.sig`
  - `terraform-provider-circleci-server_1.0.0_manifest.json`

### Provider not showing in registry
- Repository must be public
- Provider must be published through the registry UI first
- Check Terraform Registry logs for errors

## Future Releases

After initial setup, releasing is simple:

```bash
# Update version and changelog
git commit -m "Prepare v1.0.1 release"
git tag v1.0.1
git push origin v1.0.1

# GitHub Actions handles the rest automatically!
```

## Support

- Registry Documentation: https://developer.hashicorp.com/terraform/registry/providers/publishing
- GoReleaser Documentation: https://goreleaser.com/
- GitHub Actions: Check the Actions tab for build logs
