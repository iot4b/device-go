#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o iot4b_package/opt/iot4b/iot4b main.go
#GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o iot4b_package/usr/bin/iot4b main.go

cd iot4b_package 
rm -f iot4b.ipk
ls -lh opt/iot4b/iot4b | awk '{print $5}'
tar -czvf control.tar.gz CONTROL
tar -czvf data.tar.gz opt etc usr
echo 2.0 > debian-binary
/opt/homebrew/opt/binutils/bin/gar r  iot4b.ipk debian-binary control.tar.gz data.tar.gz
echo "iot4b.ipk created"
ls -lh iot4b.ipk | awk '{print $5}'
/opt/homebrew/opt/binutils/bin/gar t iot4b.ipk
echo "Cleaning up..."
rm -f control.tar.gz
rm -f data.tar.gz
rm -f debian-binary
rm -f opt/iot4b/iot4b
