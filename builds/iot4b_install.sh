#!/bin/sh

PKG_URL="https://raw.githubusercontent.com/ever-iot/docs/main/packages"
OPENWRT_PKG="iot4b_openwrt.ipk"
KEENETIC_PKG="iot4b_keenetic.ipk"
LIBNDM_PKG="libndm_1.8.0-1_mipsel-3.4_kn.ipk"
NDMQ_PKG="ndmq_1.0.2-7_mipsel-3.4_kn.ipk"

if [ -f /etc/openwrt_release ]; then
    echo "Detected OpenWRT."
    echo "Downloading package..."
    curl -o /tmp/$OPENWRT_PKG "$PKG_URL/$OPENWRT_PKG" || { echo "Failed to download $OPENWRT_PKG"; exit 1; }
    echo "Installing package..."
    opkg install /tmp/$OPENWRT_PKG
    echo "Cleaning up..."
    rm -f /tmp/$OPENWRT_PKG
    echo "Done"
elif grep -qi "NDMS" /proc/version; then
    echo "Detected Keenetic."
    echo "Downloading packages..."
    curl -o /tmp/$LIBNDM_PKG "$PKG_URL/$LIBNDM_PKG" || { echo "Failed to download $LIBNDM_PKG"; exit 1; }
    curl -o /tmp/$NDMQ_PKG "$PKG_URL/$NDMQ_PKG" || { echo "Failed to download $NDMQ_PKG"; exit 1; }
    curl -o /tmp/$KEENETIC_PKG "$PKG_URL/$KEENETIC_PKG" || { echo "Failed to download $KEENETIC_PKG"; exit 1; }
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
