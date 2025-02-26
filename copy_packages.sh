#!/bin/sh

cp  ./builds/iot4b_openwrt.ipk /opt/iot4b/repo/mipsel/iot4b.ipk

cd /opt/iot4b/repo/mipsel
opkg-make-index -a -p ./mipsel > Packages
gzip -k Packages


cp  ./builds/iot4b_keenetic.ipk /opt/iot4b/repo/mipsel-3.4_kn/iot4b.ipk
cd /opt/iot4b/repo/mipsel-3.4_kn
opkg-make-index -a -p ./mipsel-3.4_kn > Packages
gzip -k Packages