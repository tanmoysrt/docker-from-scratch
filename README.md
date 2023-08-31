# Bridge setup
ip link add name br0 type bridge
ip addr add 172.12.0.11/24 brd + dev br0
ip link set br0 up

# allow forwarding by the bridge in the root net namespace
iptables -A FORWARD -i br0 -j ACCEPT

# Expose traffice to br0 to public interface `wlo1`
iptables -t nat -A POSTROUTING -s 172.12.0.0/24 -j MASQUERADE
### br0 -> wlo1
iptables -A FORWARD -i br0 -o wlo1 -j ACCEPT
### wlo1 -> br0 for established connection
iptables -A FORWARD -i wlo1 -o br0 -m state --state RELATED,ESTABLISHED -j ACCEPT

# Setup veth pair
ip link add veth0 type veth peer name ceth0 netns <PID>
ip link set veth0 up
ip link set veth0 master br0

# Inside container configure the network
ip link set lo up
ip link set ceth0 up
ip addr add 172.12.0.12/24 dev ceth0
>> Setup gateway 
ip route add default via 172.12.0.11


# DNS setup
put in `       `
```
nameserver 8.8.8.8
```


#### Dump
```
export PATH=$PATH:/usr/local/go/bin && go run main.go run /bin/sh
```