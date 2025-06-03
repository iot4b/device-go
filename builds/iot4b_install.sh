#!/bin/sh

PKG_URL="http://repo.iot4b.co/packages/"
OPENWRT_PKG="mipsel/iot4b-mipsel.ipk"
KEENETIC_PKG="mipsel-3.4_kn/iot4b-mipsel-3.4_kn.ipk"
LIBNDM_PKG="mipsel-3.4_kn/libndm_1.8.0-1_mipsel-3.4_kn.ipk"
NDMQ_PKG="mipsel-3.4_kn/ndmq_1.0.2-7_mipsel-3.4_kn.ipk"

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
    download_and_install "$OPENWRT_PKG"
    echo "Done"
elif grep -qi "NDMS" /proc/version; then
    echo "Detected Keenetic."
    download_and_install "$LIBNDM_PKG"
    download_and_install "$NDMQ_PKG"
    download_and_install "$KEENETIC_PKG"
    echo "Done"
fi
