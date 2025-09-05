#!/bin/sh

PKG_URL="http://repo.iot4b.co/packages/"
LIBNDM_PKG_MIPS="mipsel-3.4_kn/libndm_1.8.0-1_mipsel-3.4_kn.ipk"
LIBNDM_PKG_ARM64="aarch64-3.10_kn/libndm_1.1.25-1_aarch64-3.10_kn.ipk"
NDMQ_PKG_MIPS="mipsel-3.4_kn/ndmq_1.0.2-7_mipsel-3.4_kn.ipk"
NDMQ_PKG_ARM64="aarch64-3.10_kn/ndmq_1.0.2-11_aarch64-3.10_kn.ipk"

ARCH=$(opkg print-architecture | awk '{print $2}' | tail -n1)

download_and_install() {
    local pkg=$1
    local url="$PKG_URL$pkg"
    local fname="/tmp/$(basename $pkg)"

    echo "Downloading $pkg..."
    curl -fsSL -o "$fname" "$url" || { echo "Failed to download $url"; exit 1; }

    echo "Installing $pkg..."
    opkg install "$fname"

    echo "Cleaning up $fname..."
    rm -f "$fname"
}

if [ -f /etc/openwrt_release ]; then
    echo "Detected OpenWRT."
    download_and_install "$ARCH/iot4b-openwrt-$ARCH.ipk"
    echo "Done"
elif grep -qi "NDMS" /proc/version; then
    echo "Detected Keenetic."
    if [ "$ARCH" = "mipsel-3.4_kn" ]; then
      download_and_install "$LIBNDM_PKG_MIPS"
      download_and_install "$NDMQ_PKG_MIPS"
    elif [ "$ARCH" = "aarch64-3.10_kn" ]; then
      download_and_install "$LIBNDM_PKG_ARM64"
      download_and_install "$NDMQ_PKG_ARM64"
    fi
    download_and_install "$ARCH/iot4b-keenetic-$ARCH.ipk"
    echo "Done"
fi
