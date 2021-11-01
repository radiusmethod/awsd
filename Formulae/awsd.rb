class Awsd < Formula
  desc "AWS Profile Switcher in Go"
  homepage "https://github.com/pjaudiomv/awsd"
  url "https://github.com/pjaudiomv/awsd/archive/v0.0.1.tar.gz"
  sha256 "7271f86e7e766a83b18ab107238cd6c7b28914ae"
  license "MIT"
  head "https://github.com/pjaudiomv/awsd.git", branch: "main"

  def install
    system "bash", "install.sh"
  end
end
