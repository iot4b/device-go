#!/bin/sh
set -eu

./opkg/build.sh keenetic mipsel-3.4_kn
./opkg/build.sh keenetic aarch64-3.10_kn
./opkg/build.sh openwrt mips_siflower
./opkg/build.sh openwrt armv7l
./opkg/build.sh openwrt aarch64

./apt/build.sh
