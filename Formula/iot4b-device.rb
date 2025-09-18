class Iot4bDevice < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.1.0.tar.gz"
  sha256 "5141e376f0d1faa078ab1e3930c2ab7df924318b41c3e900f191bec00c4ed5bd"
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
