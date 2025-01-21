#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o ./builds/iot4b_keenetic/opt/iot4b/iot4b main.go

cd builds/iot4b_keenetic
rm -f ../iot4b_keenetic.ipk
ls -lh opt/iot4b/iot4b | awk '{print $5}'
tar -czvf control.tar.gz control postinst
tar -czvf data.tar.gz opt tmp
echo 2.0 > debian-binary
tar -czvf ../iot4b_keenetic.ipk debian-binary control.tar.gz data.tar.gz
echo "../iot4b_keenetic.ipk created"
ls -lh ../iot4b_keenetic.ipk | awk '{print $5}'
chmod +x ../iot4b_keenetic.ipk
echo "Cleaning up..."
rm -f control.tar.gz
rm -f data.tar.gz
rm -f debian-binary
rm -f opt/iot4b/iot4b