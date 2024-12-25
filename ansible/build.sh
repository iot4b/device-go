#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o iot4b_package/opt/iot4b/iot4b main.go
ls -lh iot4b_package/opt/iot4b
cd iot4b_package
tar -czvf control.tar.gz CONTROL
tar -czvf data.tar.gz opt etc
echo 2.0 > debian-binary
tar r iot4b.ipk debian-binary control.tar.gz data.tar.gz
