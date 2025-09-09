class Iot4bDevice < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "https://github.com/iot4b/device-go/archive/refs/tags/v1.1.0.tar.gz"
  sha256 "897229a83ca74cdb8a32f529968a217b94fe88c5978de871c732b400d21c4571"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w", output: bin/"iot4b-device")
  end

  test do
    assert_match "iot4b-device version 1.1.0", shell_output("#{bin}/iot4b-device --version")
  end
end
