#!/bin/sh


rm -f /opt/iot4b/repo/packages/mipsel/*
rm -f /opt/iot4b/repo/packages/mipsel-3.4_kn/*

cp  ./builds/iot4b_openwrt.ipk /opt/iot4b/repo/packages/mipsel/iot4b-mipsel.ipk
cp  ./builds/iot4b_keenetic.ipk /opt/iot4b/repo/packages/mipsel-3.4_kn/iot4b-mipsel-3.4_kn.ipk

cd /opt/iot4b/repo/packages/mipsel
opkg-make-index -a -p ./ > Packages
gzip -k Packages


cd /opt/iot4b/repo/packages/mipsel-3.4_kn
opkg-make-index -a -p ./ > Packages
gzip -k Packages