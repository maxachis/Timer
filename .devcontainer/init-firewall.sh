#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

echo "=== Initializing firewall ==="

# Domains Claude Code needs access to
ALLOWED_DOMAINS=(
    api.anthropic.com
    statsig.anthropic.com
    statsig.com
    sentry.io
    registry.npmjs.org
    github.com
    api.github.com
    raw.githubusercontent.com
    objects.githubusercontent.com
    crates.io
    static.crates.io
    index.crates.io
    marketplace.visualstudio.com
    vscode.blob.core.windows.net
    update.code.visualstudio.com
)

# --- Phase 1: Resolve all IPs while network is still open ---

echo "Resolving domain IPs (before firewall lockdown)..."

# Create ipset for allowed domains
ipset create allowed-domains hash:net -exist
ipset flush allowed-domains

# Resolve each domain and add IPs to the ipset
for domain in "${ALLOWED_DOMAINS[@]}"; do
    echo "Resolving $domain..."
    ips=$(dig +short A "$domain" 2>/dev/null | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$' || true)
    for ip in $ips; do
        if [[ "$ip" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
            ipset add allowed-domains "$ip/32" -exist
            echo "  Added $ip ($domain)"
        fi
    done
done

# Fetch and add GitHub IP ranges
echo "Fetching GitHub IP ranges..."
gh_ranges=$(curl -s https://api.github.com/meta 2>/dev/null || true)
if echo "$gh_ranges" | jq -e '.web and .api and .git' >/dev/null 2>&1; then
    for cidr in $(echo "$gh_ranges" | jq -r '(.web + .api + .git)[]' | grep -v ':' | aggregate -q 2>/dev/null || echo "$gh_ranges" | jq -r '(.web + .api + .git)[]' | grep -v ':'); do
        if [[ "$cidr" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/[0-9]{1,2}$ ]]; then
            ipset add allowed-domains "$cidr" -exist
        fi
    done
    echo "  Added GitHub IP ranges"
else
    echo "  WARNING: Could not fetch GitHub IP ranges"
fi

# --- Phase 2: Lock down the firewall ---

echo ""
echo "Locking down firewall..."

# Preserve Docker DNS NAT rules before flushing
DOCKER_DNS_RULES=$(iptables-save -t nat | grep "127\.0\.0\.11" || true)

# Flush existing rules
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X

# Restore Docker DNS rules
if [ -n "$DOCKER_DNS_RULES" ]; then
    echo "$DOCKER_DNS_RULES" | while read -r rule; do
        iptables -t nat $rule 2>/dev/null || true
    done
fi

# Set default policies to DROP
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP

# Allow loopback
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Allow established and related connections
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

# Allow DNS (needed to resolve domains)
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT
iptables -A OUTPUT -p tcp --dport 53 -j ACCEPT

# Allow SSH
iptables -A OUTPUT -p tcp --dport 22 -j ACCEPT

# Allow host network (for Docker communication)
HOST_IP=$(ip route | grep default | cut -d" " -f3)
HOST_NETWORK=$(echo "$HOST_IP" | sed "s/\.[0-9]*$/.0\/24/")
iptables -A OUTPUT -d "$HOST_NETWORK" -j ACCEPT
iptables -A INPUT -s "$HOST_NETWORK" -j ACCEPT

# Allow traffic to IPs in the allowed-domains ipset
iptables -A OUTPUT -m set --match-set allowed-domains dst -j ACCEPT

# Reject everything else with an immediate error (not silent drop)
iptables -A OUTPUT -j REJECT --reject-with icmp-admin-prohibited

echo ""
echo "=== Firewall verification ==="

# Verify blocked traffic
if curl --connect-timeout 5 -s https://example.com >/dev/null 2>&1; then
    echo "FAIL: example.com should be blocked but is reachable"
    exit 1
else
    echo "PASS: example.com is blocked"
fi

# Verify allowed traffic
if curl --connect-timeout 5 -s https://api.anthropic.com >/dev/null 2>&1; then
    echo "PASS: api.anthropic.com is reachable"
else
    echo "WARN: api.anthropic.com not reachable (IPs may have changed, try rebuilding)"
fi

echo ""
echo "=== Firewall initialized successfully ==="
