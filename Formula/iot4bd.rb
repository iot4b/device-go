class Iot4bd < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.2.5.tar.gz"
  sha256 "640418038f2d901fdd00f72c7778be382d61b5e31bcaeff439d32a45b184094b"
  license "MIT"

  depends_on "go" => :build

  def install
    # Build binary
    system "go", "build", *std_go_args(ldflags: "-s -w", output: bin/"iot4bd")

    # Install default config
    (etc/"iot4bd").install "config/prod.yml"
  end

  service do
    run [opt_bin/"iot4bd"]
    keep_alive true
    working_dir var
    log_path var/"log/iot4bd.log"
    error_log_path var/"log/iot4bd.log"
  end

  test do
    assert_match "iot4bd version 1.2.5", shell_output("#{bin}/iot4bd --version")
  end
end
