#!/usr/bin/env sh
set -eu

if [ "$#" -ne 3 ]; then
  echo "usage: $0 <version> <url> <sha256>" >&2
  exit 1
fi

version="$1"
url="$2"
sha256="$3"

cat <<EOF
class Iot4b < Formula
  desc "IOT4B Device"
  homepage "https://github.com/iot4b/device-go"
  url "${url}"
  sha256 "${sha256}"
  license "MIT"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s
      -w
      -X device-go/packages/buildinfo.Version=${version}
    ]

    system "go", "build", *std_go_args(ldflags: ldflags, output: bin/"iot4b")

    (etc/"iot4b").install "config/iot4b.yml"
  end

  service do
    run [opt_bin/"iot4b"]
    keep_alive true
    working_dir var
    log_path var/"log/iot4b.log"
    error_log_path var/"log/iot4b.log"
  end

  test do
    assert_match "iot4b version ${version}", shell_output("#{bin}/iot4b --version")
  end
end
EOF
