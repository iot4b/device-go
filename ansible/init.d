#!/bin/sh /etc/rc.common
START=99
STOP=10

start() {
    /opt/iot4b/device.bin -env openwrt -port 5686 &
}

stop() {
    killall device.bin
}