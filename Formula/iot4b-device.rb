class Iot4bDevice < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.1.0.tar.gz"
  sha256 "8a8087b3db8669e41262f6384c81a7b6387424acb0f30e199ecaf5e6c9f18af2"
  license "MIT"

  depends_on "go" => :build

  def install
    # Build binary
    system "go", "build", *std_go_args(ldflags: "-s -w", output: bin/"iot4b-device")

    # Install default config
    (etc/"iot4b-device").install "config/prod.yaml"
  end

  test do
    assert_match "iot4b-device version 1.1.0", shell_output("#{bin}/iot4b-device --version")
  end
end
