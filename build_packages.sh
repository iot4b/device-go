#!/bin/sh

./builds/build.sh keenetic mipsel-3.4_kn
./builds/build.sh keenetic aarch64-3.10_kn
./builds/build.sh openwrt mips_siflower
./builds/build.sh openwrt armv7l
./builds/build.sh openwrt aarch64