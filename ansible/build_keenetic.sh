#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o iot4b_keenetic/opt/iot4b/iot4b main.go
#GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o iot4b_keenetic/usr/bin/iot4b main.go

cd iot4b_keenetic
rm -f iot4b.ipk
ls -lh opt/iot4b/iot4b | awk '{print $5}'
tar -czvf control.tar.gz control postinst
tar -czvf data.tar.gz opt
echo 2.0 > debian-binary
tar -czvf iot4b.ipk debian-binary control.tar.gz data.tar.gz
echo "iot4b.ipk created"
ls -lh iot4b.ipk | awk '{print $5}'
echo "Cleaning up..."
rm -f control.tar.gz
rm -f data.tar.gz
rm -f debian-binary
