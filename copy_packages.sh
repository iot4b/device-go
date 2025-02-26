#!/bin/sh

cp  ./builds/*openwrt.ipk /opt/iot4b/repo/mipsel

cd /opt/iot4b/repo/mipsel
opkg-make-index -a -p ./mipsel > Packages
gzip -k Packages


cp  ./builds/*openwrt.ipk /opt/iot4b/repo/mipsel-3.4_kn

cd /opt/iot4b/repo/mipsel-3.4_kn
opkg-make-index -a -p ./mipsel-3.4_kn > Packages
gzip -k Packages