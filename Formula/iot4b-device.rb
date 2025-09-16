class Iot4bDevice < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.1.0.tar.gz"
  sha256 "dab0b074c69aac8ca703308b05db6583e66d306336af6fc5e9feb44a9e675ec6"
  license "MIT"

  depends_on "go" => :build

  def install
    # Build binary
    system "go", "build", *std_go_args(ldflags: "-s -w", output: bin/"iot4b-device")

    # Install default config
    (etc/"iot4b-device").install "config/prod.yml"
  end

  test do
    assert_match "iot4b-device version 1.1.0", shell_output("#{bin}/iot4b-device --version")
  end
end
