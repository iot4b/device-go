class Iot4b < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/1.2.7.tar.gz"
  sha256 "c1901b5ccb5b1436c9236f060dd427672b1fa3df370ef7f939c56984ab69db07"
  license "MIT"

  depends_on "go" => :build

  def install
    # Build binary
    system "go", "build", *std_go_args(ldflags: "-s -w", output: bin/"iot4b")

    # Install default config
    (etc/"iot4b").install "config/prod.yml"
  end

  service do
    run [opt_bin/"iot4b"]
    keep_alive true
    working_dir var
    log_path var/"log/iot4b.log"
    error_log_path var/"log/iot4b.log"
  end

  test do
    assert_match "iot4b version 1.2.7", shell_output("#{bin}/iot4b --version")
  end
end
