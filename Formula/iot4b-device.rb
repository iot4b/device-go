class Iot4bDevice < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.2.3.tar.gz"
  sha256 "5e34ac830b7b20ef6f719eb0b4b1599e671b90b7937af9cedb426c0e28da53af"
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
    assert_match "iot4b-device version 1.2.3", shell_output("#{bin}/iot4b-device --version")
  end
end
