# 🚀 Quick Publishing Guide

Follow these steps to publish `sportsterminal` to GitHub and Homebrew.

## Prerequisites

- GitHub account (elliota43)
- Git installed locally
- GitHub personal access token (for goreleaser)

## Method 1: Automated with GoReleaser (Recommended)

This method automatically creates releases and updates Homebrew.

### Step 1: Install GoReleaser

```bash
brew install goreleaser
```

### Step 2: Create GitHub Repositories

1. **Main repository**: Create `sportsterminal` at https://github.com/new
   - Repository name: `sportsterminal`
   - Public repository
   - Don't initialize with README (we have one)

2. **Homebrew tap**: Create `homebrew-tap` at https://github.com/new
   - Repository name: `homebrew-tap`
   - Public repository
   - Don't initialize with anything

### Step 3: Push Your Code

```bash
cd /Users/elliotanderson/sports-scores

# Initialize and push
git init
git add .
git commit -m "Initial commit: sportsterminal v1.0.0"
git branch -M main
git remote add origin https://github.com/elliota43/sportsterminal.git
git push -u origin main
```

### Step 4: Get GitHub Token

1. Go to https://github.com/settings/tokens/new
2. Note: "GoReleaser"
3. Expiration: Your choice
4. Select scopes:
   - `repo` (Full control of private repositories)
   - `write:packages` (Upload packages)
5. Click "Generate token"
6. Copy the token

```bash
export GITHUB_TOKEN="your_token_here"
```

### Step 5: Create Release

```bash
# Create and push tag
git tag -a v1.0.0 -m "First release"
git push origin v1.0.0

# Run goreleaser
goreleaser release --clean
```

**Done!** 🎉 GoReleaser will:
- Build binaries for all platforms
- Create GitHub release with binaries
- Automatically create/update your Homebrew tap

### Step 6: Test Installation

```bash
brew tap elliota43/tap
brew install sportsterminal
sportsterminal --version
```

## Method 2: Manual (Without GoReleaser)

If you prefer manual control:

### Step 1: Push to GitHub

```bash
cd /Users/elliotanderson/sports-scores
git init
git add .
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/elliota43/sportsterminal.git
git push -u origin main
```

### Step 2: Create GitHub Release

1. Go to https://github.com/elliota43/sportsterminal/releases/new
2. Tag version: `v1.0.0`
3. Release title: `v1.0.0 - Initial Release`
4. Description: Add your release notes
5. Click "Publish release"

### Step 3: Calculate SHA256

```bash
curl -sL https://github.com/elliota43/sportsterminal/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
```

Copy the hash output.

### Step 4: Set Up Homebrew Tap

```bash
# Clone your homebrew tap
git clone https://github.com/elliota43/homebrew-tap.git
cd homebrew-tap

# Create Formula directory
mkdir -p Formula

# Copy the formula
cp /Users/elliotanderson/sports-scores/Formula/sportsterminal.rb Formula/

# Edit the formula and add the SHA256 from Step 3
# Replace the empty sha256 "" with the actual hash
nano Formula/sportsterminal.rb
```

Update line 5:
```ruby
sha256 "your_sha256_hash_here"
```

### Step 5: Push Formula

```bash
git add Formula/sportsterminal.rb
git commit -m "Add sportsterminal formula v1.0.0"
git push origin main
```

### Step 6: Test

```bash
brew tap elliota43/tap
brew install sportsterminal
sportsterminal
```

## Future Updates

### With GoReleaser:

```bash
# Update version in main.go if needed
# Commit your changes
git add .
git commit -m "Update for v1.1.0"
git push

# Create new release
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0
goreleaser release --clean
```

### Manual Updates:

1. Create new GitHub release with new tag
2. Get new SHA256
3. Update formula in homebrew-tap
4. Push changes

## Troubleshooting

**GoReleaser fails with auth error:**
- Check your GITHUB_TOKEN is valid
- Ensure token has `repo` and `write:packages` scopes
- Try: `export GITHUB_TOKEN="your_token"`

**Homebrew tap not found:**
- Ensure repository is public
- Verify repository name is exactly `homebrew-tap`
- Check repo exists at `https://github.com/elliota43/homebrew-tap`

**Formula fails to install:**
- Verify SHA256 matches the release tarball
- Test build locally: `go build .`
- Check formula syntax

**Users can't find your tap:**
```bash
# They need to run:
brew tap elliota43/tap
brew install sportsterminal
```

## What Users Will Do

Once published, users install with:

```bash
# One-time: Add your tap
brew tap elliota43/tap

# Install
brew install sportsterminal

# Or in one command:
brew install elliota43/tap/sportsterminal
```

## Files Created for Publishing

- ✅ `Formula/sportsterminal.rb` - Homebrew formula
- ✅ `.goreleaser.yml` - GoReleaser configuration  
- ✅ `HOMEBREW.md` - Detailed Homebrew guide
- ✅ `scripts/publish.sh` - Helper script (manual method)
- ✅ `main.go` - Added version flag support

## Next Steps

1. Choose Method 1 (GoReleaser) or Method 2 (Manual)
2. Follow the steps above
3. Test the installation
4. Share with users: `brew install elliota43/tap/sportsterminal`

---

**Need help?** See `HOMEBREW.md` for more detailed information.

