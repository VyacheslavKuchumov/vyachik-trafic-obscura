#!/bin/sh
ip addr add 10.0.0.1/24 dev tun0
ip link set tun0 up

sysctl -w net.ipv4.ip_forward=1
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE
