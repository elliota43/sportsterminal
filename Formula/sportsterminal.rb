class Sportsterminal < Formula
  desc "Beautiful terminal interface for checking live sports scores"
  homepage "https://github.com/elliota43/sportsterminal"
  url "https://github.com/elliota43/sportsterminal/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "" # This will be calculated after creating the release
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"sportsterminal"
  end

  test do
    # Test that the binary exists and is executable
    assert_match "Error running program", shell_output("#{bin}/sportsterminal 2>&1", 1)
  end
end

