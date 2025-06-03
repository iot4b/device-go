#!/bin/sh

PKG_URL="http://repo.iot4b.co/packages/"
OPENWRT_PKG="ihttp://repo.iot4b.co/packages/mipsel/iot4b-mipsel.ipk"
KEENETIC_PKG="http://repo.iot4b.co/packages/mipsel-3.4_kn/iot4b-mipsel-3.4_kn.ipk"
LIBNDM_PKG="http://repo.iot4b.co/packages/mipsel-3.4_kn/libndm_1.8.0-1_mipsel-3.4_kn.ipk"
NDMQ_PKG="http://repo.iot4b.co/packages/mipsel-3.4_kn/ndmq_1.0.2-7_mipsel-3.4_kn.ipk"

if [ -f /etc/openwrt_release ]; then
    echo "Detected OpenWRT."
    echo "Downloading package..."
    curl -o /tmp/$OPENWRT_PKG "$OPENWRT_PKG" || { echo "Failed to download $OPENWRT_PKG"; exit 1; }
    echo "Installing package..."
    opkg install /tmp/$OPENWRT_PKG
    echo "Cleaning up..."
    rm -f /tmp/$OPENWRT_PKG
    echo "Done"
elif grep -qi "NDMS" /proc/version; then
    echo "Detected Keenetic."
    echo "Downloading packages..."
    curl -o /tmp/$LIBNDM_PKG "$LIBNDM_PKG" || { echo "Failed to download $LIBNDM_PKG"; exit 1; }
    curl -o /tmp/$NDMQ_PKG "$NDMQ_PKG" || { echo "Failed to download $NDMQ_PKG"; exit 1; }
    curl -o /tmp/$KEENETIC_PKG "$KEENETIC_PKG" || { echo "Failed to download $KEENETIC_PKG"; exit 1; }
    echo "Installing packages..."
    opkg install /tmp/$LIBNDM_PKG
    opkg install /tmp/$NDMQ_PKG
    opkg install /tmp/$KEENETIC_PKG
    echo "Cleaning up..."
    rm -f /tmp/$LIBNDM_PKG
    rm -f /tmp/$NDMQ_PKG
    rm -f /tmp/$KEENETIC_PKG
    echo "Done"
fi
