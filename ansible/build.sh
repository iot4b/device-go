#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o device.bin main.go
ls -lh device.bin
