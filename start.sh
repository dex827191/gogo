#!/bin/bash

ulimit -n 65535 2>/dev/null
ulimit -s unlimited 2>/dev/null

SYS_CFG="/tmp/system_config.json"

if [ -f "$SYS_CFG" ] && command -v jq >/dev/null 2>&1; then
    CUSTOM_BANNER=$(jq -r '.banner // empty' "$SYS_CFG")
    ENABLE_BBR=$(jq -r '.enable_bbr // "true"' "$SYS_CFG")
    UDPGW_PORT=$(jq -r '.udpgw_port // "7300"' "$SYS_CFG")
    UDPGW_MAX_CLIENTS=$(jq -r '.udpgw_max_clients // "1000"' "$SYS_CFG")
else
    CUSTOM_BANNER=""
    ENABLE_BBR="true"
    UDPGW_PORT="7300"
    UDPGW_MAX_CLIENTS="1000"
fi

echo "[*] Kernel tweaks..."
sysctl -w net.ipv4.tcp_tw_reuse=1 2>/dev/null
sysctl -w net.ipv4.tcp_fin_timeout=10 2>/dev/null
sysctl -w net.core.default_qdisc=fq 2>/dev/null
sysctl -w net.ipv4.tcp_rmem="4096 8388608 16777216" 2>/dev/null
sysctl -w net.ipv4.tcp_wmem="4096 8388608 16777216" 2>/dev/null
sysctl -w net.core.rmem_max=16777216 2>/dev/null
sysctl -w net.core.wmem_max=16777216 2>/dev/null
sysctl -w net.core.netdev_max_backlog=50000 2>/dev/null
sysctl -w net.ipv4.tcp_max_syn_backlog=8192 2>/dev/null

if [ "$ENABLE_BBR" = "true" ]; then
    sysctl -w net.ipv4.tcp_congestion_control=bbr 2>/dev/null
else
    sysctl -w net.ipv4.tcp_congestion_control=cubic 2>/dev/null
fi

# Banner
if [ -n "$CUSTOM_BANNER" ]; then
    echo -e "$CUSTOM_BANNER" > /etc/dropbear_banner
else
    cat << 'EOF' > /etc/dropbear_banner
==================================================
          👑 SELAMAT MENIKMATI 👑
       🥳 SSH SERVER PAAS RAILWAY 🥳
==================================================
 🔹 MULTIPLEXER : GO ENGINE
 🔹 OS PLATFORM : UBUNTU
 🔹 SSH SERVICE : DROPBEAR ENHANCED BUFFER
==================================================
 powered by : d e d e f a t h u
==================================================
EOF
fi

# Stunnel
SSL_INTERNAL_PORT="2443"
mkdir -p /etc/stunnel /var/run/stunnel4
openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 \
    -subj "/C=ID/ST=Jakarta/L=Jakarta/O=RailwaySSH/CN=localhost" \
    -keyout /etc/stunnel/stunnel.pem -out /etc/stunnel/stunnel.pem 2>/dev/null
chmod 600 /etc/stunnel/stunnel.pem

pkill -9 dropbear 2>/dev/null
/usr/sbin/dropbear -p 127.0.0.1:22 -b /etc/dropbear_banner -W 1048576 -K 15 -I 300
sleep 1

cat <<EOF > /etc/stunnel/stunnel.conf
pid = /var/run/stunnel4/stunnel.pid
foreground = no
debug = 0

[ssh-ssl]
accept = 127.0.0.1:$SSL_INTERNAL_PORT
connect = 127.0.0.1:22
cert = /etc/stunnel/stunnel.pem
EOF

rm -f /var/run/stunnel4/stunnel.pid 2>/dev/null
stunnel4 /etc/stunnel/stunnel.conf

# UDPGW
if [ -f /usr/local/bin/badvpn-udpgw ]; then
    echo "[*] Starting badvpn-udpgw on ${UDPGW_PORT}..."
    /usr/local/bin/badvpn-udpgw --listen-addr 0.0.0.0:${UDPGW_PORT} \
        --max-clients ${UDPGW_MAX_CLIENTS} --max-connections-for-client 50 &
fi

# Named tunnel via env (opsional)
if [ -n "$TOKEN" ]; then
    echo "[*] Starting Zero Trust tunnel..."
    /usr/local/bin/cloudflared tunnel run --protocol http2 --no-tls-verify \
        --token "$TOKEN" --url "http://localhost:${ARGO_PORT:-8001}" \
        > /tmp/named_tunnel.log 2>&1 &
fi

# Go binary handle: mux (8881) + wsproxy (8880) + server (8081)
exec /app/panel
