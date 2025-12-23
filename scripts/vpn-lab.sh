#!/usr/bin/env bash
set -e

CLIENT_NS="vpn-client"
SERVER_NS="vpn-server"

CLIENT_TUN="tun1"
SERVER_TUN="tun0"

CLIENT_VETH="veth-client"
SERVER_VETH="veth-server"

VPN_NET="10.0.0.0/24"
CLIENT_VPN_IP="10.0.0.2/24"
SERVER_VPN_IP="10.0.0.1/24"

LINK_NET="192.168.100.0/24"
CLIENT_LINK_IP="192.168.100.2/24"
SERVER_LINK_IP="192.168.100.1/24"

cleanup() {
    echo "[*] Cleaning up"
    ip netns del $CLIENT_NS 2>/dev/null || true
    ip netns del $SERVER_NS 2>/dev/null || true
    ip link del $CLIENT_VETH 2>/dev/null || true
    ip tuntap del $CLIENT_TUN mode tun 2>/dev/null || true
    ip tuntap del $SERVER_TUN mode tun 2>/dev/null || true
    iptables -t nat -D POSTROUTING -s $VPN_NET -j MASQUERADE 2>/dev/null || true
}

if [[ "$1" == "clean" ]]; then
    cleanup
    exit 0
fi

cleanup

echo "[*] Creating network namespaces"
ip netns add $SERVER_NS
ip netns add $CLIENT_NS

echo "[*] Enabling loopback"
ip netns exec $SERVER_NS ip link set lo up
ip netns exec $CLIENT_NS ip link set lo up

echo "[*] Creating TUN devices"
ip tuntap add $SERVER_TUN mode tun
ip tuntap add $CLIENT_TUN mode tun

ip link set $SERVER_TUN netns $SERVER_NS
ip link set $CLIENT_TUN netns $CLIENT_NS


echo "[*] Assigning VPN IP addresses"
ip netns exec $SERVER_NS ip addr add $SERVER_VPN_IP dev $SERVER_TUN
ip netns exec $CLIENT_NS ip addr add $CLIENT_VPN_IP dev $CLIENT_TUN

ip netns exec $SERVER_NS ip link set $SERVER_TUN up
ip netns exec $CLIENT_NS ip link set $CLIENT_TUN up

echo "[*] Creating veth pair (simulated internet)"
ip link add $CLIENT_VETH type veth peer name $SERVER_VETH

ip link set $CLIENT_VETH netns $CLIENT_NS
ip link set $SERVER_VETH netns $SERVER_NS

ip netns exec $CLIENT_NS ip addr add $CLIENT_LINK_IP dev $CLIENT_VETH
ip netns exec $SERVER_NS ip addr add $SERVER_LINK_IP dev $SERVER_VETH

ip netns exec $CLIENT_NS ip link set $CLIENT_VETH up
ip netns exec $SERVER_NS ip link set $SERVER_VETH up

echo "[*] Enabling IP forwarding in server namespace"
ip netns exec $SERVER_NS sysctl -w net.ipv4.ip_forward=1 >/dev/null

echo "[*] Enabling NAT on server"
iptables -t nat -A POSTROUTING -s $VPN_NET -j MASQUERADE

echo
echo "[✓] VPN test environment is ready"
echo
echo "Run server:"
echo "  sudo ip netns exec $SERVER_NS ./bin/server"
echo
echo "Run client:"
echo "  sudo ip netns exec $CLIENT_NS ./bin/client"
echo
echo "Test tunnel:"
echo "  sudo ip netns exec $CLIENT_NS ping 10.0.0.1"
echo
echo "Cleanup:"
echo "  sudo ./vpn-lab.sh clean"
