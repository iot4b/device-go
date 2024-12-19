#!/bin/sh

GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o package/iot4b/files/iot4b main.go
ls -lh package/iot4b/iot4b
