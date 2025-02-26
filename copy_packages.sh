#!/bin/sh



cp  ./builds/iot4b_openwrt.ipk /opt/iot4b/repo/mipsel/iot4b-mipsel.ipk
cp  ./builds/iot4b_keenetic.ipk /opt/iot4b/repo/mipsel-3.4_kn/iot4b-mipsel-3.4_kn.ipk


cd /opt/iot4b/repo/mipsel
rm -f *
opkg-make-index -a -p ./mipsel > Packages
gzip -k Packages


cd /opt/iot4b/repo/mipsel-3.4_kn
rm -f *
opkg-make-index -a -p ./mipsel-3.4_kn > Packages
gzip -k Packages