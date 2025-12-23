#!/bin/sh
ip addr add 10.0.0.2/24 dev tun0
ip link set tun0 up
ip route add default dev tun0
