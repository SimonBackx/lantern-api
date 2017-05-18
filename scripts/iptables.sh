
echo "flush iptable rules"
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X
iptables -t mangle -F
iptables -t mangle -X

echo "Set default policy to 'DROP'"
iptables -P INPUT   DROP
iptables -P FORWARD DROP
iptables -P OUTPUT  DROP

## This should be one of the first rules.
## so dns lookups are already allowed for your other rules

echo "Allowing DNS lookups (tcp, udp port 53)"
iptables -A OUTPUT -p udp --dport 53 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT  -p udp --sport 53 -m state --state ESTABLISHED     -j ACCEPT
iptables -A OUTPUT -p tcp --dport 53 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT  -p tcp --sport 53 -m state --state ESTABLISHED     -j ACCEPT

echo "allow all and everything on localhost"
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

#######################################################################################################
## Global iptable rules. Not IP specific

echo "Allowing new and established outgoing connections to port 80, 443"
iptables -A OUTPUT  -p tcp --dport 80 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp --sport 80 -m state --state ESTABLISHED     -j ACCEPT

iptables -A OUTPUT  -p tcp --dport 443 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp --sport 443 -m state --state ESTABLISHED     -j ACCEPT

echo "Allowing incoming connections to 443"
iptables -A INPUT  -p tcp --dport 443 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A OUTPUT -p tcp --sport 443 -m state --state ESTABLISHED     -j ACCEPT

echo "Allow all outgoing connections to port 22"
iptables -A OUTPUT -p tcp --dport 22 -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT  -p tcp --sport 22 -m state --state ESTABLISHED     -j ACCEPT

echo "Allow incoming ssh"
iptables -A INPUT   -p tcp   --sport 513:65535   --dport 22        -m state --state NEW,ESTABLISHED -j ACCEPT
iptables -A OUTPUT  -p tcp   --sport 22          --dport 513:65535 -m state --state ESTABLISHED     -j ACCEPT

echo "Set default policy to 'DROP' for ipv6"
ip6tables -P INPUT   DROP
ip6tables -P FORWARD DROP
ip6tables -P OUTPUT  DROP

echo "Saving rules for next reboot"
sudo su -c 'iptables-save > /etc/iptables/rules.v4'
sudo su -c 'ip6tables-save > /etc/iptables/rules.v6'

exit 0