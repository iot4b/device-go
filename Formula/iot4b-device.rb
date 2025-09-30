class Iot4bDevice < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.2.2.tar.gz"
  sha256 "c148e184c0825f59b55c4b08b77342e4b904ca789593754b6c4522349917b4a1"
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
    assert_match "iot4b-device version 1.2.2", shell_output("#{bin}/iot4b-device --version")
  end
end
