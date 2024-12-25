#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o iot4b_package/opt/iot4b/iot4b main.go
cd iot4b_package 
ls -lh opt/iot4b/iot4b | awk '{print $5}'
tar -czvf control.tar.gz CONTROL
tar -czvf data.tar.gz opt etc
echo 2.0 > debian-binary
tar -cf iot4b.ipk debian-binary control.tar.gz data.tar.gz
echo "iot4b.ipk created"
ls -lh iot4b.ipk | awk '{print $5}'
echo "Cleaning up..."
rm -f control.tar.gz
rm -f data.tar.gz
rm -f debian-binary
rm -f opt/iot4b/iot4b
