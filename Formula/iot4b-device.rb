class Iot4bDevice < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.2.0.tar.gz"
  sha256 "089cb2776e400515e84be576f5210206f46b80465b4387dfbe01c5292585b0f8"
  license "MIT"

  depends_on "go" => :build

  def install
    # Build binary
    system "go", "build", *std_go_args(ldflags: "-s -w", output: bin/"iot4b-device")

    # Install default config
    (etc/"iot4b-device").install "config/prod.yml"
  end

  service do
    run [opt_bin/"iot4b-device"]
    keep_alive true
    working_dir var
    log_path var/"log/iot4b-device.log"
    error_log_path var/"log/iot4b-device.log"
  end

  test do
    assert_match "iot4b-device version 1.1.0", shell_output("#{bin}/iot4b-device --version")
  end
end
