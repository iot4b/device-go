#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o ./builds/iot4b_openwrt/opt/iot4b/iot4b main.go

cd builds/iot4b_openwrt
rm -f ../iot4b_openwrt.ipk
ls -lh opt/iot4b/iot4b | awk '{print $5}'
tar -czvf control.tar.gz control
tar -czvf data.tar.gz opt etc
echo 2.0 > debian-binary
#/opt/homebrew/opt/binutils/bin/gar r ../iot4b_openwrt.ipk debian-binary control.tar.gz data.tar.gz
tar -czvf ../iot4b_openwrt.ipk debian-binary control.tar.gz data.tar.gz
echo "Created ./builds/iot4b_openwrt.ipk, size: $(ls -lh ../iot4b_openwrt.ipk | awk '{print $5}')"
echo "Cleaning up..."
rm -f control.tar.gz  data.tar.gz debian-binary opt/iot4b/iot4b