# 🍺 Publishing to Homebrew

This guide will help you publish `sportsterminal` to Homebrew so users can install it with `brew install sportsterminal`.

## Option 1: Homebrew Tap (Recommended for Personal Projects)

This is the easiest way to distribute your app via Homebrew.

### Step 1: Create a Homebrew Tap Repository

1. Create a new GitHub repository named `homebrew-tap`:
   ```bash
   # On GitHub, create: https://github.com/elliota43/homebrew-tap
   ```

2. Clone and set it up:
   ```bash
   git clone https://github.com/elliota43/homebrew-tap.git
   cd homebrew-tap
   mkdir -p Formula
   ```

### Step 2: Push Your Code and Create a Release

1. First, push your sportsterminal code to GitHub:
   ```bash
   cd /path/to/sportsterminal
   git init
   git add .
   git commit -m "Initial commit"
   git branch -M main
   git remote add origin https://github.com/elliota43/sportsterminal.git
   git push -u origin main
   ```

2. Create a GitHub release (v1.0.0):
   - Go to https://github.com/elliota43/sportsterminal/releases
   - Click "Create a new release"
   - Tag: `v1.0.0`
   - Title: `v1.0.0 - Initial Release`
   - Description: Add release notes
   - Click "Publish release"

3. Get the SHA256 of your release tarball:
   ```bash
   curl -sL https://github.com/elliota43/sportsterminal/archive/refs/tags/v1.0.0.tar.gz | shasum -a 256
   ```

### Step 3: Update and Publish the Formula

1. Copy the formula to your tap:
   ```bash
   cp Formula/sportsterminal.rb /path/to/homebrew-tap/Formula/
   ```

2. Update the `sha256` field in the formula with the hash from Step 2.3

3. Commit and push:
   ```bash
   cd /path/to/homebrew-tap
   git add Formula/sportsterminal.rb
   git commit -m "Add sportsterminal formula"
   git push origin main
   ```

### Step 4: Users Can Now Install!

Users can install your app with:
```bash
brew tap elliota43/tap
brew install sportsterminal
```

Or in one line:
```bash
brew install elliota43/tap/sportsterminal
```

## Option 2: Homebrew Core (For Popular Projects)

To get into the official Homebrew repository:

1. Your project should be notable and have significant usage
2. Must be stable (version 1.0.0+)
3. Submit a PR to [Homebrew/homebrew-core](https://github.com/Homebrew/homebrew-core)
4. Follow their [contribution guidelines](https://docs.brew.sh/Adding-Software-to-Homebrew)

## Testing Your Formula Locally

Before publishing, test the formula:

```bash
# Install from local formula
brew install --build-from-source Formula/sportsterminal.rb

# Test it works
sportsterminal

# Uninstall
brew uninstall sportsterminal
```

## Updating Your Formula for New Releases

When you release a new version:

1. Create a new GitHub release with a new tag (e.g., v1.1.0)
2. Get the new SHA256:
   ```bash
   curl -sL https://github.com/elliota43/sportsterminal/archive/refs/tags/v1.1.0.tar.gz | shasum -a 256
   ```
3. Update the formula:
   - Change the `url` version
   - Update the `sha256`
4. Commit and push to your tap

## Auto-Update with goreleaser (Optional)

You can automate releases with [GoReleaser](https://goreleaser.com/):

1. Install goreleaser:
   ```bash
   brew install goreleaser
   ```

2. Create `.goreleaser.yml` in your project root (see the file created below)

3. Create releases automatically:
   ```bash
   git tag -a v1.0.0 -m "First release"
   git push origin v1.0.0
   goreleaser release --clean
   ```

This will automatically:
- Build for multiple platforms
- Create GitHub releases
- Update your Homebrew tap

## Quick Reference

### Installation Commands Users Will Use:

```bash
# Add your tap
brew tap elliota43/tap

# Install
brew install sportsterminal

# Upgrade
brew upgrade sportsterminal

# Uninstall
brew uninstall sportsterminal
```

## README Updates

Add this to your README.md installation section:

```markdown
### Homebrew (macOS/Linux)

\`\`\`bash
brew install elliota43/tap/sportsterminal
\`\`\`
```

## Troubleshooting

**Formula doesn't work?**
- Verify the SHA256 matches the release tarball
- Check the URL is correct and accessible
- Test locally first

**Build fails?**
- Ensure go.mod is properly configured
- Test `go build` works in the project root
- Check Go version compatibility

---

For more info, see the [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook).

