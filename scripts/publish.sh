#!/bin/bash
set -e

# Script to help publish sportsterminal to GitHub and set up Homebrew tap

VERSION="v1.0.0"
REPO="elliota43/sportsterminal"
TAP_REPO="elliota43/homebrew-tap"

echo "🚀 Publishing sportsterminal to GitHub and Homebrew"
echo "=================================================="
echo ""

# Check if we're in a git repo
if [ ! -d .git ]; then
    echo "📦 Initializing git repository..."
    git init
    git add .
    git commit -m "Initial commit: sportsterminal v1.0.0"
    git branch -M main
    echo "✅ Git repository initialized"
    echo ""
fi

echo "📋 Next steps to publish:"
echo ""
echo "1️⃣  Create GitHub repository:"
echo "   Go to: https://github.com/new"
echo "   Repository name: sportsterminal"
echo "   Make it public"
echo ""
echo "2️⃣  Push your code:"
echo "   git remote add origin https://github.com/${REPO}.git"
echo "   git push -u origin main"
echo ""
echo "3️⃣  Create a release:"
echo "   Go to: https://github.com/${REPO}/releases/new"
echo "   Tag version: ${VERSION}"
echo "   Release title: ${VERSION} - Initial Release"
echo "   Click 'Publish release'"
echo ""
echo "4️⃣  Get the SHA256 hash:"
echo "   Run this command after creating the release:"
echo "   curl -sL https://github.com/${REPO}/archive/refs/tags/${VERSION}.tar.gz | shasum -a 256"
echo ""
echo "5️⃣  Create Homebrew tap repository:"
echo "   Go to: https://github.com/new"
echo "   Repository name: homebrew-tap"
echo "   Make it public"
echo "   Clone it: git clone https://github.com/${TAP_REPO}.git"
echo ""
echo "6️⃣  Set up the tap:"
echo "   cd homebrew-tap"
echo "   mkdir -p Formula"
echo "   cp /path/to/sportsterminal/Formula/sportsterminal.rb Formula/"
echo "   # Update the sha256 field with the hash from step 4"
echo "   git add Formula/sportsterminal.rb"
echo "   git commit -m 'Add sportsterminal formula'"
echo "   git push origin main"
echo ""
echo "7️⃣  Users can now install with:"
echo "   brew install ${TAP_REPO}/sportsterminal"
echo ""
echo "📚 For detailed instructions, see HOMEBREW.md"
echo ""
echo "🎉 Or use GoReleaser for automated releases (recommended):"
echo "   brew install goreleaser"
echo "   export GITHUB_TOKEN=your_github_token"
echo "   git tag -a ${VERSION} -m 'First release'"
echo "   git push origin ${VERSION}"
echo "   goreleaser release --clean"

