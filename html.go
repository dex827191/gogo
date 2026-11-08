package main

// HTML UI di-embed langsung ke binary — tidak butuh file eksternal.
// Menggunakan raw string literal (backtick), jadi semua karakter tersimpan literal.
const indexHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <title>SSH & VPN PANEL</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: '-apple-system', BlinkMacSystemFont, sans-serif; background: #000000; color: #ffffff; min-height: 100vh; padding: 0; margin: 0; }
        .navbar { background: #000000; border-bottom: 1px solid #262626; padding: 10px 15px; display: flex; align-items: center; justify-content: space-between; position: sticky; top: 0; z-index: 100; }
        .nav-left { display: flex; align-items: center; gap: 12px; }
        .hamburger-btn { background: none; border: 1px solid #262626; color: #ffffff; font-size: 18px; cursor: pointer; display: flex; align-items: center; justify-content: center; width: 32px; height: 32px; border-radius: 6px; }
        .hamburger-btn:hover { background: #1a1a1a; }
        .nav-title { font-size: 13px; font-weight: bold; color: #ffffff; text-transform: uppercase; letter-spacing: 0.5px; }
        .sidebar-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); z-index: 998; opacity: 0; pointer-events: none; transition: opacity 0.3s ease; }
        .sidebar-overlay.active { opacity: 1; pointer-events: auto; }
        .sidebar { position: fixed; top: 0; left: -280px; width: 280px; height: 100%; background: #000000; border-right: 1px solid #262626; z-index: 999; transition: left 0.3s ease; display: flex; flex-direction: column; padding: 20px 0; }
        .sidebar.active { left: 0; }
        .sidebar-header { padding: 0 20px 20px 20px; border-bottom: 1px solid #262626; display: flex; justify-content: space-between; align-items: center; }
        .sidebar-title { font-size: 13px; font-weight: bold; color: #ffffff; }
        .close-sidebar { background: none; border: none; color: #ffffff; font-size: 18px; cursor: pointer; }
        .sidebar-menu { list-style: none; padding: 15px 0; margin: 0; flex: 1; overflow-y: auto; }
        .menu-item { padding: 12px 20px; display: flex; align-items: center; gap: 12px; color: #ffffff; font-size: 13px; font-weight: 600; cursor: pointer; transition: 0.2s; border-left: 3px solid transparent; }
        .menu-item:hover, .menu-item.active { background: #141414; color: #ffffff; border-left-color: #ffffff; }
        .main-wrapper { display: flex; justify-content: center; padding: 15px 12px; }
        .container { background: #000000; width: 100%; max-width: 480px; }
        .page-section { display: none; }
        .page-section.active { display: block; }
        .dev-tag { font-size: 10px; color: #ffffff; margin-top: 2px; font-weight: bold; text-align: center; }
        .btn-login-trigger { background: #000000; color: #ffffff; border: 1px solid #262626; padding: 5px 10px; border-radius: 6px; font-size: 10px; cursor: pointer; font-weight: bold; }
        .btn-login-trigger:hover { background: #1a1a1a; }
        .status-container { text-align: center; margin-bottom: 12px; }
        .status-badge { display: inline-block; background: #000000; color: #ffffff; padding: 4px 10px; border-radius: 50px; font-size: 10px; font-weight: bold; border: 1px solid #262626; }
        .stats-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-bottom: 15px; }
        .stat-card { background: #000000; padding: 8px 10px; border-radius: 6px; border: 1px solid #262626; text-align: left; }
        .stat-title { font-size: 9px; color: #ffffff; text-transform: uppercase; font-weight: bold; }
        .stat-value { font-size: 12px; font-weight: bold; color: #ffffff; margin-top: 2px; }
        .ssh-manager { background: #000000; padding: 12px; border-radius: 10px; border: 1px solid #262626; margin-bottom: 15px; }
        .ssh-title { font-size: 12px; font-weight: bold; color: #ffffff; text-transform: uppercase; margin-bottom: 8px; display: flex; justify-content: space-between; align-items: center; }
        .input-group { display: flex; gap: 6px; margin-bottom: 8px; }
        .input-ssh { background: #000000; border: 1px solid #262626; padding: 6px 10px; border-radius: 6px; color: #ffffff; font-size: 12px; width: 100%; outline: none; }
        .input-ssh:focus { border-color: #ffffff; }
        .btn-add { background: #000000; color: #ffffff; border: 1px solid #262626; padding: 6px 12px; border-radius: 6px; font-weight: bold; cursor: pointer; font-size: 11px; }
        .btn-add:hover { background: #1a1a1a; }
        .admin-status-lbl { font-size: 9px; font-weight: bold; color: #ffffff; border: 1px solid #262626; padding: 2px 5px; border-radius: 4px; }
        .result-box { display: none; background: #000000; border: 1px solid #262626; border-radius: 6px; padding: 8px; font-family: monospace; font-size: 11px; color: #ffffff; white-wrap: pre-wrap; margin-bottom: 10px; overflow-x: hidden; }
        .btn-copy-result { display: none; background: #000000; color: #ffffff; border: 1px solid #262626; padding: 5px; border-radius: 4px; font-weight: bold; cursor: pointer; font-size: 10px; width: 100%; margin-bottom: 10px; }
        .btn-copy-result:hover { background: #1a1a1a; }
        .ssh-list { width: 100%; border-collapse: collapse; margin-top: 8px; font-size: 11px; }
        .ssh-list th { text-align: left; padding: 4px; color: #ffffff; border-bottom: 1px solid #262626; }
        .ssh-list td { padding: 4px; border-bottom: 1px solid #262626; vertical-align: middle; color: #ffffff; }
        .btn-action-group { display: flex; gap: 4px; justify-content: flex-end; }
        .btn-del { background: #000000; color: #ffffff; border: 1px solid #262626; padding: 3px 6px; border-radius: 4px; cursor: pointer; font-size: 10px; display: none; }
        .btn-del:hover { background: #1a1a1a; }
        .btn-info { background: #000000; color: #ffffff; border: 1px solid #262626; padding: 3px 6px; border-radius: 4px; cursor: pointer; font-size: 10px; font-weight: bold; display: none; }
        .btn-info:hover { background: #1a1a1a; }
        .url-section { background: #000000; border: 1px solid #262626; padding: 8px 10px; border-radius: 8px; margin-bottom: 8px; }
        .url-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px; }
        .url-title { font-size: 10px; font-weight: bold; text-transform: uppercase; letter-spacing: 0.3px; color: #ffffff; }
        .url-content-row { display: flex; align-items: center; gap: 8px; background: #000000; border: 1px solid #262626; padding: 4px 8px; border-radius: 6px; }
        .url-box { font-family: monospace; font-size: 11px; word-break: break-all; font-weight: bold; flex: 1; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; color: #ffffff; }
        .btn-copy-mini { background: #000000; color: #ffffff; border: 1px solid #262626; padding: 4px 8px; border-radius: 4px; font-size: 10px; font-weight: bold; cursor: pointer; white-space: nowrap; transition: 0.2s; }
        .btn-copy-mini:hover { background: #1a1a1a; }
        .select-zt-mini { background: #000000; border: none; color: #ffffff; font-size: 11px; font-weight: bold; font-family: monospace; outline: none; width: 100%; cursor: pointer; }
        .btn-token-trigger { width: 100%; padding: 12px; border-radius: 8px; font-weight: bold; font-size: 12px; cursor: pointer; border: 1px solid #262626; background: #000000; color: #ffffff; transition: 0.2s; text-transform: uppercase; letter-spacing: 0.5px; }
        .btn-token-trigger:hover { background: #1a1a1a; }
        .card-blue { background-color: #000000; border: 1px solid #262626; padding: 12px; border-radius: 10px; margin-top: 10px; text-align: left; }
        .btn-blue { background-color: #000000; border: 1px solid #262626; color: #ffffff; padding: 6px; border-radius: 6px; font-size: 10px; font-weight: bold; cursor: pointer; width: 100%; text-align: center; font-family: monospace; }
        .btn-blue:hover { background-color: #1a1a1a; }
        .btn-active { border-color: #ffffff !important; background-color: #1a1a1a !important; color: #ffffff !important; }
        .grid-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 6px; margin-top: 6px; }
        .grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 6px; margin-bottom: 8px; }
        .lbl-vpn { font-size: 9px; color: #ffffff; font-weight: bold; display: block; margin-bottom: 3px; text-transform: uppercase; }
        .border-lbl { border-left: 2px solid #ffffff; padding-left: 6px; font-size: 10px; font-weight: bold; margin-top: 10px; font-family: monospace; color: #ffffff; }
        .zt-admin-card { background: #000000; border: 1px solid #262626; padding: 12px; border-radius: 10px; margin-bottom: 12px; }
        .sub-box { background: #000000; border: 1px solid #262626; padding: 10px; border-radius: 8px; margin-bottom: 10px; }
        .sub-box-title { font-size: 11px; font-weight: bold; color: #ffffff; margin-bottom: 6px; text-transform: uppercase; display: flex; align-items: center; gap: 6px; }
        .note { font-size: 10px; color: #ffffff; text-align: center; line-height: 1.4; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="navbar">
        <div class="nav-left">
            <button class="hamburger-btn" onclick="toggleSidebar()">☰</button>
            <div class="nav-title">SERVER PANEL</div>
        </div>
        <button class="btn-login-trigger" id="admin-login-btn" onclick="handleAdminAuthBtn()">LOGIN ADMIN</button>
    </div>

    <div class="sidebar-overlay" id="sidebar-overlay" onclick="toggleSidebar()"></div>
    <div class="sidebar" id="sidebar">
        <div class="sidebar-header">
            <div class="sidebar-title">MENU NAVIGATION</div>
            <button class="close-sidebar" onclick="toggleSidebar()">✕</button>
        </div>
        <ul class="sidebar-menu">
            <li class="menu-item active" onclick="switchPage('page-dashboard')">Dashboard Utama</li>
            <li class="menu-item" onclick="switchPage('page-buat-ssh')">Buat Akun SSH</li>
            <li class="menu-item" onclick="switchPage('page-set-ssh')">Pengaturan SSH</li>
            <li class="menu-item" onclick="switchPage('page-set-v2ray')">Pengaturan V2Ray / VMess</li>
            <li class="menu-item" onclick="switchPage('page-set-cfip')">Pengaturan CFIP</li>
            <li class="menu-item" onclick="switchPage('page-generator')">Config Generator</li>
        </ul>
        <div style="padding: 15px 20px; border-top: 1px solid #262626;">
            <span id="btn-change-pass" onclick="changeAdminPassUI()" style="color: #ffffff; cursor: pointer; font-size: 11px; text-decoration: underline; display: none;">GANTI PASS ADMIN</span>
        </div>
    </div>

    <div class="main-wrapper">
        <div class="container">
            <div id="page-dashboard" class="page-section active">
                <div style="text-align:center; margin-bottom: 10px;">
                    <h1 style="font-size:16px; color:#ffffff; text-transform:uppercase;">DASHBOARD SERVER</h1>
                    <div class="dev-tag">DYNAMIC TRIPLE-TUNNEL NODE CORE ACTIVE</div>
                </div>
                <div class="status-container"><div class="status-badge"><span style="color: #ffffff">ALL TUNNELS ONLINE</span></div></div>
                <div class="stats-grid">
                    <div class="stat-card" style="grid-column: span 2;"><div class="stat-title">CPU Model</div><div class="stat-value" id="cpu" style="font-size:11px; color:#ffffff;">Loading...</div></div>
                    <div class="stat-card"><div class="stat-title">RAM Used / Total</div><div class="stat-value" id="ram">Loading...</div></div>
                    <div class="stat-card"><div class="stat-title">Disk Usage (/)</div><div class="stat-value" id="disk">Loading...</div></div>
                    <div class="stat-card"><div class="stat-title">Server Uptime</div><div class="stat-value" id="uptime" style="font-size:11px;">Loading...</div></div>
                    <div class="stat-card"><div class="stat-title">SSH Online</div><div class="stat-value" id="ssh" style="font-size:12px; line-height:1.2;">0 Users</div></div>
                </div>
                <div style="margin-bottom: 12px;">
                    <button class="btn-token-trigger" onclick="promptSingleTokenInput()">MASUKKAN TOKEN ARGO TUNNEL</button>
                </div>
                <div class="sub-box" style="margin-bottom: 12px;">
                    <div class="sub-box-title" style="text-align:center; display:block;">INFO PENGATURAN SERVER</div>
                    <div style="font-size: 10px; color: #ffffff; text-align: center; font-weight: bold; line-height: 1.5;">
                        <span>CFIP: </span><span id="display-cfip" style="font-family:monospace;">Loading...</span> |
                        <span>DNS: </span><span id="display-dns" style="font-family:monospace;">UDP</span> |
                        <span>ENGINE: </span><span id="display-engine" style="font-family:monospace;">WS</span><br>
                        <span>PORT: </span><span id="display-ws-port" style="font-family:monospace;">22</span> |
                        <span>KEEPALIVE: </span><span id="display-ws-keep" style="font-family:monospace;">15000ms</span> |
                        <span>BBR: </span><span id="display-bbr" style="font-family:monospace;">ON</span> |
                        <span>UDPGW: </span><span id="display-udpgw" style="font-family:monospace;">7300</span>
                    </div>
                </div>
                <div class="url-section">
                    <div class="url-header"><span class="url-title">SSH WS Domain</span></div>
                    <div class="url-content-row">
                        <div id="zt-container" style="flex:1; overflow:hidden;"><div class="url-box" id="named-url">Menghubungkan...</div></div>
                        <button class="btn-copy-mini" id="btn-copy-named" onclick="copyTxt('named-url', 'btn-copy-named')">Copy</button>
                    </div>
                </div>
                <div class="url-section">
                    <div class="url-header"><span class="url-title">VMess / VLess Domain</span></div>
                    <div class="url-content-row">
                        <div id="zt-vmess-container" style="flex:1; overflow:hidden;"><div class="url-box" id="vmess-named-url">Menghubungkan...</div></div>
                        <button class="btn-copy-mini" id="btn-copy-vmess-named" onclick="copyTxt('vmess-named-url', 'btn-copy-vmess-named')">Copy</button>
                    </div>
                </div>
                <div class="url-section">
                    <div class="url-header"><span class="url-title">SSH SNI Server</span></div>
                    <div class="url-content-row">
                        <div class="url-box" id="railway-url">Loading...</div>
                        <button class="btn-copy-mini" id="btn-copy-railway" onclick="copyTxt('railway-url', 'btn-copy-railway')">Copy</button>
                    </div>
                </div>
                <div class="url-section">
                    <div class="url-header"><span class="url-title">Quick Tunnel Sub</span></div>
                    <div class="url-content-row">
                        <div class="url-box" id="quick-url">Loading...</div>
                        <button class="btn-copy-mini" id="btn-copy-quick" onclick="copyTxt('quick-url', 'btn-copy-quick')">Copy</button>
                    </div>
                </div>
            </div>

            <div id="page-buat-ssh" class="page-section">
                <div class="ssh-manager">
                    <div class="ssh-title"><span>BUAT AKUN SSH BARU</span><span id="admin-indicator" class="admin-status-lbl">PUBLIC CREATION</span></div>
                    <div class="input-group">
                        <input type="text" id="ssh-user" class="input-ssh" placeholder="Username...">
                        <input type="password" id="ssh-pass" class="input-ssh" placeholder="Password...">
                        <button class="btn-add" id="btn-add-ssh" onclick="createAccount()">ADD</button>
                    </div>
                    <div id="ssh-result" class="result-box"></div>
                    <button id="btn-copy-acc" class="btn-copy-result" onclick="copyAccountText()">COPY DETAIL AKUN</button>
                    <div id="ssh-msg" style="font-size: 11px; margin-top: 5px; font-weight: bold; color: #ffffff;"></div>
                    <div class="ssh-title" style="margin-top: 15px; border-top: 1px solid #262626; padding-top: 10px;">DAFTAR AKUN TERDAFTAR</div>
                    <table class="ssh-list">
                        <thead><tr><th>Username</th><th>Shell Path</th><th style="text-align: right;">Aksi</th></tr></thead>
                        <tbody id="ssh-table-body"><tr><td colspan="3" style="text-align:center; color:#ffffff;">Loading accounts...</td></tr></tbody>
                    </table>
                </div>
            </div>

            <div id="page-set-ssh" class="page-section">
                <div class="zt-admin-card">
                    <div class="sub-box" style="margin-bottom:0;">
                        <div class="sub-box-title">PENGATURAN SSH SERVER</div>
                        <div class="grid-3" style="margin-top:4px;">
                            <div>
                                <label class="lbl-vpn">PORT DROPBEAR</label>
                                <select id="wsPortDropdown" class="input-ssh" onchange="toggleCustomWsPortInput()">
                                    <option value="22">Port 22 (OpenSSH)</option>
                                    <option value="109">Port 109 (Dropbear Direct)</option>
                                    <option value="143">Port 143 (Dropbear Alt)</option>
                                    <option value="447">Port 447 (SSL Stunnel)</option>
                                    <option value="custom">Custom...</option>
                                </select>
                                <input type="number" id="wsPortInput" class="input-ssh" placeholder="Port..." style="display:none; margin-top:4px;">
                            </div>
                            <div>
                                <label class="lbl-vpn">KEEPALIVE</label>
                                <select id="wsKeepDropdown" class="input-ssh" onchange="toggleCustomWsKeepInput()">
                                    <option value="5000">5000 (5s)</option>
                                    <option value="15000">15000 (15s)</option>
                                    <option value="30000">30000 (30s)</option>
                                    <option value="custom">Custom...</option>
                                </select>
                                <input type="number" id="wsKeepInput" class="input-ssh" placeholder="MS..." style="display:none; margin-top:4px;">
                            </div>
                            <div>
                                <label class="lbl-vpn">MAX BUFFER</label>
                                <select id="wsBufDropdown" class="input-ssh" onchange="toggleCustomWsBufInput()">
                                    <option value="16384">16384 (16KB)</option>
                                    <option value="32768">32768 (32KB)</option>
                                    <option value="65536">65536 (64KB)</option>
                                    <option value="custom">Custom...</option>
                                </select>
                                <input type="number" id="wsBufInput" class="input-ssh" placeholder="Bytes..." style="display:none; margin-top:4px;">
                            </div>
                        </div>
                        <div style="margin-top:10px; border-top:1px solid #262626; padding-top:8px;">
                            <div>
                                <label class="lbl-vpn">BANNER DROPBEAR</label>
                                <textarea id="bannerInput" class="input-ssh" style="height:50px; font-family:monospace; font-size:10px;" placeholder="Kosongkan untuk banner standar..."></textarea>
                            </div>
                            <div class="grid-2" style="margin-top:6px; margin-bottom:0;">
                                <div>
                                    <label class="lbl-vpn">TCP BBR</label>
                                    <select id="bbrSelect" class="input-ssh" style="font-weight:bold;">
                                        <option value="true">ON</option>
                                        <option value="false">OFF</option>
                                    </select>
                                </div>
                                <div>
                                    <label class="lbl-vpn">UDPGW PORT</label>
                                    <select id="udpgwPortSelect" class="input-ssh" style="font-weight:bold;">
                                        <option value="7300">Port 7300</option>
                                        <option value="7200">Port 7200</option>
                                        <option value="7100">Port 7100</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <button class="btn-token-trigger" style="margin-top:8px; padding:8px;" onclick="saveWsProxySettingUI()">SIMPAN SSH CONFIG</button>
                    </div>
                </div>
            </div>

            <div id="page-set-v2ray" class="page-section">
                <div class="zt-admin-card">
                    <div class="sub-box" style="margin-bottom:0;">
                        <div class="sub-box-title">PENGATURAN V2RAY SERVER</div>
                        <div style="display: flex; flex-direction: column; gap: 6px;">
                            <div class="grid-2" style="margin-bottom:0;">
                                <div>
                                    <label class="lbl-vpn">DNS MODE</label>
                                    <select id="dnsTypeSelect" class="input-ssh" style="font-weight: bold;" onchange="toggleDnsPlaceholder()">
                                        <option value="udp">UDP Fast DNS</option>
                                        <option value="doh">DoH Secure DNS</option>
                                    </select>
                                </div>
                                <div>
                                    <label class="lbl-vpn">ENGINE MODE</label>
                                    <select id="engineSelect" class="input-ssh" style="font-weight: bold;">
                                        <option value="ws">WebSocket (WS)</option>
                                        <option value="grpc">gRPC Engine</option>
                                        <option value="h2">HTTP/2 (H2)</option>
                                        <option value="tcp">TCP Direct</option>
                                    </select>
                                </div>
                            </div>
                            <div>
                                <label class="lbl-vpn">CUSTOM DNS IP / DOH URL</label>
                                <select id="dnsDropdown" class="input-ssh" onchange="toggleCustomDnsInput()">
                                    <option value="8.8.8.8">Google DNS (8.8.8.8)</option>
                                    <option value="1.1.1.1">Cloudflare DNS (1.1.1.1)</option>
                                    <option value="9.9.9.9">Quad9 DNS (9.9.9.9)</option>
                                    <option value="https://1.1.1.1/dns-query">Cloudflare DoH</option>
                                    <option value="https://dns.google/dns-query">Google DoH</option>
                                    <option value="custom">Custom...</option>
                                </select>
                                <input type="text" id="customDnsInput" class="input-ssh" placeholder="Ketik DNS Manual..." style="display:none; margin-top:4px;">
                            </div>
                            <button class="btn-token-trigger" style="margin-top: 4px; padding:8px;" onclick="saveDnsNetworkSetting()">SIMPAN V2RAY CONFIG</button>
                        </div>
                    </div>
                </div>
            </div>

            <div id="page-set-cfip" class="page-section">
                <div class="zt-admin-card">
                    <div class="sub-box" style="margin-bottom:0;">
                        <div class="sub-box-title">CLOUDFLARE CLEAN IP (CFIP)</div>
                        <p style="font-size: 10px; color: #ffffff; margin-bottom: 10px; line-height: 1.4;">Atur IP Cloudflare yang bersih untuk kestabilan koneksi.</p>
                        <button class="btn-token-trigger" style="padding:10px;" onclick="promptCfipInput()">UBAH / RESET CFIP</button>
                    </div>
                </div>
            </div>

            <div id="page-generator" class="page-section">
                <div class="card-blue" style="margin-top:0;">
                    <div style="text-align: center; margin-bottom: 10px; border-bottom: 1px solid #262626; padding-bottom: 6px;">
                        <span style="font-size: 12px; font-weight: bold; color: #ffffff;">CONFIG GENERATOR</span>
                    </div>
                    <div class="grid-2">
                        <div>
                            <label class="lbl-vpn">UUID / PASS</label>
                            <input id="uuidInput" type="text" value="Loading..." class="input-ssh" style="font-family: monospace;" readonly>
                        </div>
                        <div>
                            <label class="lbl-vpn">TARGET DOMAIN</label>
                            <select id="domainSelect" class="input-ssh" style="font-family: monospace; font-weight: bold;">
                                <option value="">-- Menunggu Domain --</option>
                            </select>
                        </div>
                    </div>
                    <div style="margin-bottom: 10px;">
                        <label class="lbl-vpn">BUG HOST (SNI / CDN)</label>
                        <input id="bugInput" type="text" value="suporte.garena.com" class="input-ssh" style="font-family: monospace;">
                    </div>
                    <div class="border-lbl">BUG SNI (NORMAL)</div>
                    <div class="grid-3">
                        <button onclick="buildConfig('vless', 'sni', event)" class="btn-blue">VLESS</button>
                        <button onclick="buildConfig('vmess', 'sni', event)" class="btn-blue">VMESS</button>
                        <button onclick="buildConfig('trojan', 'sni', event)" class="btn-blue">TROJAN</button>
                    </div>
                    <div class="border-lbl">BUG SNI (REVERSE)</div>
                    <div class="grid-3">
                        <button onclick="buildConfig('vless', 'sni_reverse', event)" class="btn-blue">VLESS REV</button>
                        <button onclick="buildConfig('vmess', 'sni_reverse', event)" class="btn-blue">VMESS REV</button>
                        <button onclick="buildConfig('trojan', 'sni_reverse', event)" class="btn-blue">TROJAN REV</button>
                    </div>
                    <div class="border-lbl">BUG CDN (PROXY)</div>
                    <div class="grid-3">
                        <button onclick="buildConfig('vless', 'cdn', event)" class="btn-blue">VLESS</button>
                        <button onclick="buildConfig('vmess', 'cdn', event)" class="btn-blue">VMESS</button>
                        <button onclick="buildConfig('trojan', 'cdn', event)" class="btn-blue">TROJAN</button>
                    </div>
                    <div id="output-area" class="result-box" style="margin-top: 10px; display: none;">
                        <div style="display: flex; justify-content: space-between; font-weight: bold; margin-bottom: 4px; font-size: 10px;">
                            <span id="out-type">CONFIG</span>
                            <span onclick="copyOutConfig()" style="cursor: pointer; text-decoration: underline;">[COPY]</span>
                        </div>
                        <p id="configText" style="word-break: break-all; color: #ffffff; max-height: 80px; overflow-y: auto; font-family: monospace;"></p>
                    </div>
                </div>
            </div>

            <p class="note">Single terowongan multi-host berjalan sinkron.<br>Go Core Engine Rendering System.</p>
        </div>
    </div>

    <script>
        let adminToken = localStorage.getItem("admin_session_token") || "";
        let savedUsersData = [];
        let isPassConfigured = false;

        function toggleSidebar() {
            document.getElementById('sidebar').classList.toggle('active');
            document.getElementById('sidebar-overlay').classList.toggle('active');
        }

        function switchPage(pageId) {
            document.querySelectorAll('.page-section').forEach(page => {
                page.classList.remove('active');
            });
            document.getElementById(pageId).classList.add('active');
            document.querySelectorAll('.menu-item').forEach(item => {
                item.classList.remove('active');
            });
            if (event && event.currentTarget) {
                event.currentTarget.classList.add('active');
            }
            toggleSidebar();
        }

        function checkAdminUI() {
            let indicator = document.getElementById('admin-indicator');
            let loginBtn = document.getElementById('admin-login-btn');
            let changePassBtn = document.getElementById('btn-change-pass');
            if(!isPassConfigured) {
                loginBtn.innerText = "SETUP PASS";
                loginBtn.style.background = "#000000"; loginBtn.style.color = "#ffffff";
                if(changePassBtn) changePassBtn.style.display = "none";
            } else if(adminToken) {
                if(indicator) { indicator.innerText = "ADMIN ROUTE"; indicator.style.color = "#ffffff"; indicator.style.background = "#000000"; }
                loginBtn.innerText = "LOGOUT"; loginBtn.style.background = "#000000"; loginBtn.style.color = "#ffffff";
                if(changePassBtn) changePassBtn.style.display = "inline";
                document.querySelectorAll('.btn-del').forEach(b => b.style.display = "inline-block"); document.querySelectorAll('.btn-info').forEach(b => b.style.display = "inline-block");
            } else {
                if(indicator) { indicator.innerText = "PUBLIC CREATION"; indicator.style.color = "#ffffff"; indicator.style.background = "#000000"; }
                loginBtn.innerText = "LOGIN ADMIN"; loginBtn.style.background = "#000000"; loginBtn.style.color = "#ffffff";
                if(changePassBtn) changePassBtn.style.display = "none";
                document.querySelectorAll('.btn-del').forEach(b => b.style.display = "none"); document.querySelectorAll('.btn-info').forEach(b => b.style.display = "none");
            }
        }

        function toggleDnsPlaceholder() { toggleCustomDnsInput(); }
        function toggleCustomDnsInput() {
            let sel = document.getElementById('dnsDropdown').value;
            let inp = document.getElementById('customDnsInput');
            inp.style.display = (sel === 'custom') ? 'block' : 'none';
        }
        function toggleCustomWsPortInput() {
            let sel = document.getElementById('wsPortDropdown').value;
            let inp = document.getElementById('wsPortInput');
            inp.style.display = (sel === 'custom') ? 'block' : 'none';
        }
        function toggleCustomWsKeepInput() {
            let sel = document.getElementById('wsKeepDropdown').value;
            let inp = document.getElementById('wsKeepInput');
            inp.style.display = (sel === 'custom') ? 'block' : 'none';
        }
        function toggleCustomWsBufInput() {
            let sel = document.getElementById('wsBufDropdown').value;
            let inp = document.getElementById('wsBufInput');
            inp.style.display = (sel === 'custom') ? 'block' : 'none';
        }

        async function handleAdminAuthBtn() {
            if(!isPassConfigured) {
                let newP = prompt("KREASI PASSWORD ADMIN PERTAMA KALI:\nMasukkan Password Admin Baru:");
                if(!newP) return false;
                try {
                    let res = await fetch('/api/setup-pass?pass=' + encodeURIComponent(newP));
                    let data = await res.json();
                    alert(data.message);
                    if(data.status === "success") {
                        adminToken = newP;
                        localStorage.setItem("admin_session_token", adminToken);
                        updateStats();
                        return true;
                    }
                } catch(e) { alert("Gagal membuat password!"); }
                return false;
            }
            if(adminToken) {
                localStorage.removeItem("admin_session_token");
                adminToken = "";
                checkAdminUI();
                fetchAccounts();
                return false;
            }
            let pass = prompt("Masukkan Password Admin:"); if(!pass) return false;
            try {
                let res = await fetch('/api/login?pass='+pass); let data = await res.json();
                if(data.status === "success") {
                    adminToken = data.token;
                    localStorage.setItem("admin_session_token", adminToken);
                    checkAdminUI();
                    fetchAccounts();
                    return true;
                } else {
                    alert(data.message);
                }
            } catch(e) { alert("Gagal terhubung"); }
            return false;
        }

        async function saveDnsNetworkSetting() {
            if (!adminToken) {
                alert("Akses Ditolak! Anda harus Login Admin terlebih dahulu.");
                let loggedIn = await handleAdminAuthBtn();
                if (!loggedIn && !adminToken) return;
            }
            let dnsType = document.getElementById('dnsTypeSelect').value;
            let selDns = document.getElementById('dnsDropdown').value;
            let customDns = (selDns === 'custom') ? document.getElementById('customDnsInput').value.trim() : selDns;
            let engine = document.getElementById('engineSelect').value;
            if (!customDns) customDns = "8.8.8.8";
            try {
                let res = await fetch('/api/set-network?pass=' + encodeURIComponent(adminToken) + '&dns_type=' + dnsType + '&custom_dns=' + encodeURIComponent(customDns) + '&engine=' + engine);
                let data = await res.json();
                alert(data.message);
                if (data.status === "success") updateStats();
            } catch(e) { alert("Gagal menyimpan konfigurasi ke backend!"); }
        }

        async function saveWsProxySettingUI() {
            if (!adminToken) {
                alert("Akses Ditolak! Anda harus Login Admin terlebih dahulu.");
                let loggedIn = await handleAdminAuthBtn();
                if (!loggedIn && !adminToken) return;
            }
            let selPort = document.getElementById('wsPortDropdown').value;
            let port = (selPort === 'custom') ? document.getElementById('wsPortInput').value.trim() : selPort;
            let selKeep = document.getElementById('wsKeepDropdown').value;
            let keep = (selKeep === 'custom') ? document.getElementById('wsKeepInput').value.trim() : selKeep;
            let selBuf = document.getElementById('wsBufDropdown').value;
            let buf = (selBuf === 'custom') ? document.getElementById('wsBufInput').value.trim() : selBuf;
            let banner = document.getElementById('bannerInput').value.trim();
            let bbr = document.getElementById('bbrSelect').value;
            let udpgw = document.getElementById('udpgwPortSelect').value;
            if (!port) port = "22";
            if (!keep) keep = "15000";
            if (!buf) buf = "32768";
            try {
                let res1 = await fetch('/api/set-wsproxy?pass=' + encodeURIComponent(adminToken) + '&ssh_port=' + port + '&keep_alive=' + keep + '&max_buffer=' + buf);
                let data1 = await res1.json();
                let res2 = await fetch('/api/set-system?pass=' + encodeURIComponent(adminToken) + '&banner=' + encodeURIComponent(banner) + '&enable_bbr=' + bbr + '&udpgw_port=' + udpgw);
                let data2 = await res2.json();
                alert(data1.message + "\n" + data2.message);
                updateStats();
            } catch(e) { alert("Gagal memperbarui SSH Config!"); }
        }

        async function promptSingleTokenInput() {
            if (!adminToken) {
                alert("Akses Ditolak! Anda harus Login Admin terlebih dahulu.");
                let loggedIn = await handleAdminAuthBtn();
                if (!loggedIn && !adminToken) return;
            }
            let inputToken = prompt("MASUKKAN TOKEN CLOUDFLARE ARGO:\n\n(Kosongkan lalu klik OK jika ingin menghapus token tersimpan)");
            if (inputToken === null) return;
            try {
                let res = await fetch('/api/set-token?pass=' + encodeURIComponent(adminToken) + '&token=' + encodeURIComponent(inputToken.trim()));
                let data = await res.json();
                alert(data.message);
                updateStats();
            } catch(e) { alert("Gagal memperbarui token!"); }
        }

        async function promptCfipInput() {
            if (!adminToken) {
                alert("Akses Ditolak! Anda harus Login Admin terlebih dahulu.");
                let loggedIn = await handleAdminAuthBtn();
                if (!loggedIn && !adminToken) return;
            }
            let inputIp = prompt("MASUKKAN IP CLEAN CLOUDFLARE (CFIP):\nContoh: 104.17.3.81\n\n(Kosongkan untuk reset ke default)");
            if (inputIp === null) return;
            try {
                let res = await fetch('/api/set-cfip?pass=' + encodeURIComponent(adminToken) + '&ip=' + encodeURIComponent(inputIp.trim()));
                let data = await res.json();
                alert(data.message);
                updateStats();
            } catch(e) { alert("Gagal memperbarui CFIP!"); }
        }

        async function changeAdminPassUI() {
            if(!adminToken) return;
            let newP = prompt("GANTI PASSWORD ADMIN:\nMasukkan Password Admin Baru:");
            if(!newP) return;
            try {
                let res = await fetch('/api/setup-pass?old_pass=' + encodeURIComponent(adminToken) + '&pass=' + encodeURIComponent(newP));
                let data = await res.json();
                alert(data.message);
                if(data.status === "success") { adminToken = newP; localStorage.setItem("admin_session_token", adminToken); checkAdminUI(); }
            } catch(e) { alert("Gagal mengupdate password!"); }
        }

        async function updateStats() {
            try {
                let res = await fetch('/api/stats'); let data = await res.json();
                isPassConfigured = data.pass_configured;
                checkAdminUI();
                document.getElementById('cpu').innerText = data.cpu_model || "N/A";
                document.getElementById('ram').innerText = (data.ram_used || "0") + " / " + (data.ram_total || "0");
                document.getElementById('disk').innerText = data.disk_usage || "0%";
                document.getElementById('uptime').innerText = data.uptime || "0 Hours";
                let detailActiveList = data.user_list_details || "Semua user offline";
                document.getElementById('ssh').innerHTML = (data.ssh_online || "0") + " Active<br><span style='font-size:9px; font-weight:normal; color:#ffffff; display:block; margin-top:2px; white-space:pre-line;'>" + detailActiveList + "</span>";
                document.getElementById('display-cfip').innerText = data.active_cfip || "Default";
                if(data.dns_type) {
                    document.getElementById('dnsTypeSelect').value = data.dns_type;
                    document.getElementById('display-dns').innerText = (data.dns_type + " (" + (data.custom_dns || "8.8.8.8") + ")").toUpperCase();
                }
                if(data.custom_dns) {
                    let dropdown = document.getElementById('dnsDropdown');
                    let found = Array.from(dropdown.options).some(opt => opt.value === data.custom_dns);
                    if (found) {
                        dropdown.value = data.custom_dns;
                    } else {
                        dropdown.value = 'custom';
                        document.getElementById('customDnsInput').value = data.custom_dns;
                    }
                    toggleCustomDnsInput();
                }
                if(data.engine_mode) {
                    document.getElementById('engineSelect').value = data.engine_mode;
                    document.getElementById('display-engine').innerText = data.engine_mode.toUpperCase();
                }
                if(data.ws_proxy_cfg) {
                    let pDrop = document.getElementById('wsPortDropdown');
                    let foundP = Array.from(pDrop.options).some(opt => opt.value == data.ws_proxy_cfg.sshPort);
                    if (foundP) { pDrop.value = data.ws_proxy_cfg.sshPort; } else { pDrop.value = 'custom'; document.getElementById('wsPortInput').value = data.ws_proxy_cfg.sshPort; }
                    toggleCustomWsPortInput();
                    let kDrop = document.getElementById('wsKeepDropdown');
                    let foundK = Array.from(kDrop.options).some(opt => opt.value == data.ws_proxy_cfg.keepAlive);
                    if (foundK) { kDrop.value = data.ws_proxy_cfg.keepAlive; } else { kDrop.value = 'custom'; document.getElementById('wsKeepInput').value = data.ws_proxy_cfg.keepAlive; }
                    toggleCustomWsKeepInput();
                    let bDrop = document.getElementById('wsBufDropdown');
                    let foundB = Array.from(bDrop.options).some(opt => opt.value == data.ws_proxy_cfg.maxBuffer);
                    if (foundB) { bDrop.value = data.ws_proxy_cfg.maxBuffer; } else { bDrop.value = 'custom'; document.getElementById('wsBufInput').value = data.ws_proxy_cfg.maxBuffer; }
                    toggleCustomWsBufInput();
                    document.getElementById('display-ws-port').innerText = data.ws_proxy_cfg.sshPort;
                    document.getElementById('display-ws-keep').innerText = (data.ws_proxy_cfg.keepAlive || 15000) + "ms";
                }
                if(data.sys_settings) {
                    if(data.sys_settings.banner) document.getElementById('bannerInput').value = data.sys_settings.banner;
                    if(data.sys_settings.enable_bbr) document.getElementById('bbrSelect').value = data.sys_settings.enable_bbr;
                    if(data.sys_settings.udpgw_port) document.getElementById('udpgwPortSelect').value = data.sys_settings.udpgw_port;
                    document.getElementById('display-bbr').innerText = (data.sys_settings.enable_bbr === "true") ? "ON" : "OFF";
                    document.getElementById('display-udpgw').innerText = data.sys_settings.udpgw_port || "7300";
                }
                let ztContainer = document.getElementById('zt-container');
                if (data.zt_domains && data.zt_domains.length > 1) {
                    let dropdownHtml = '<select id="named-url" class="select-zt-mini" onmousedown="event.stopPropagation()">';
                    data.zt_domains.forEach(item => { dropdownHtml += '<option value="' + item.domain + '">' + item.domain + ' (' + item.port + ')</option>'; });
                    dropdownHtml += '</select>';
                    ztContainer.innerHTML = dropdownHtml;
                } else if (data.zt_domains && data.zt_domains.length === 1) {
                    ztContainer.innerHTML = '<div class="url-box" id="named-url">' + data.zt_domains[0].domain + '</div>';
                } else {
                    ztContainer.innerHTML = '<div class="url-box" id="named-url">Menghubungkan...</div>';
                }
                let ztVmessContainer = document.getElementById('zt-vmess-container');
                if (data.zt_vmess_domains && data.zt_vmess_domains.length > 1) {
                    let dropdownVmessHtml = '<select id="vmess-named-url" class="select-zt-mini" onmousedown="event.stopPropagation()">';
                    data.zt_vmess_domains.forEach(item => { dropdownVmessHtml += '<option value="' + item.domain + '">' + item.domain + ' (' + item.port + ')</option>'; });
                    dropdownVmessHtml += '</select>';
                    ztVmessContainer.innerHTML = dropdownVmessHtml;
                } else if (data.zt_vmess_domains && data.zt_vmess_domains.length === 1) {
                    ztVmessContainer.innerHTML = '<div class="url-box" id="vmess-named-url">' + data.zt_vmess_domains[0].domain + '</div>';
                } else {
                    ztVmessContainer.innerHTML = '<div class="url-box" id="vmess-named-url">Menghubungkan...</div>';
                }
                document.getElementById('railway-url').innerText = data.railway_url || "Tidak Aktif";
                document.getElementById('quick-url').innerText = data.quick_url || "Menunggu Quick Tunnel...";
                let domainSelect = document.getElementById('domainSelect');
                let currentSelected = domainSelect.value;
                let optionsHtml = '';
                if (data.quick_url && !data.quick_url.includes("Menunggu")) {
                    optionsHtml += '<option value="' + data.quick_url + '">Quick: ' + data.quick_url + '</option>';
                }
                if (data.zt_vmess_domains && data.zt_vmess_domains.length > 0) {
                    data.zt_vmess_domains.forEach(d => {
                        optionsHtml += '<option value="' + d.domain + '">Zero Argo: ' + d.domain + '</option>';
                    });
                }
                if (!optionsHtml) {
                    optionsHtml = '<option value="">-- Menunggu Domain --</option>';
                }
                domainSelect.innerHTML = optionsHtml;
                if (currentSelected) domainSelect.value = currentSelected;
            } catch(e) {}
        }

        async function fetchAccounts() {
            try {
                let res = await fetch('/api/list'); let data = await res.json(); let tbody = document.getElementById('ssh-table-body'); tbody.innerHTML = "";
                if(data.status === "success" && data.users.length > 0) {
                    savedUsersData = data.users;
                    data.users.forEach(u => {
                        tbody.innerHTML += '<tr><td style="font-weight:bold; color:#ffffff;">'+u.username+'</td><td style="color:#ffffff;">'+u.shell+'</td><td style="text-align: right;"><div class="btn-action-group"><button class="btn-info" onclick="showAccountDetails(\''+u.username+'\')">INFO</button><button class="btn-del" onclick="deleteAccount(\''+u.username+'\')">HAPUS</button></div></td></tr>';
                    });
                    checkAdminUI();
                } else { tbody.innerHTML = '<tr><td colspan="3" style="text-align:center; color:#ffffff;">Belum ada akun SSH kustom</td></tr>'; }
            } catch(e) {}
        }

        function showAccountDetails(username) {
            let userObj = savedUsersData.find(u => u.username === username);
            if(userObj) {
                alert("DATA AKUN:\n===============================\nUsername   : " + userObj.username + "\nPassword   : " + userObj.password + "\nIP Address : " + userObj.ip);
            }
        }

        async function createAccount() {
            let user = document.getElementById('ssh-user').value.trim(); let pass = document.getElementById('ssh-pass').value.trim(); let msg = document.getElementById('ssh-msg'); let resBox = document.getElementById('ssh-result'); let copyBtn = document.getElementById('btn-copy-acc');
            if(!user || !pass) { msg.style.color = "#ffffff"; msg.innerText = "Isi username & password!"; return; }
            try {
                let res = await fetch('/api/add?user='+user+'&pass='+pass); let data = await res.json();
                if(data.status === "success") { msg.innerText = ""; resBox.innerText = data.message; resBox.style.display = "block"; copyBtn.style.display = "block"; document.getElementById('ssh-user').value = ""; document.getElementById('ssh-pass').value = ""; fetchAccounts(); } else { msg.style.color = "#ffffff"; msg.innerText = data.message; resBox.style.display = "none"; copyBtn.style.display = "none"; }
            } catch(e) { msg.innerText = "Gagal memproses API"; }
        }

        function copyAccountText() {
            let txt = document.getElementById('ssh-result').innerText;
            navigator.clipboard.writeText(txt);
            let btn = document.getElementById('btn-copy-acc');
            btn.innerText = "COPIED!";
            setTimeout(() => { btn.innerText = "COPY DETAIL AKUN"; }, 1500);
        }

        async function deleteAccount(username) {
            if(!adminToken) { alert("Aksi Ilegal! Login Admin dulu Bos!"); return; }
            if(confirm("Hapus akun SSH "+username+"?")) {
                try {
                    let res = await fetch('/api/delete?user='+username+'&token='+adminToken); let data = await res.json();
                    if(data.status === "success") { fetchAccounts(); } else { alert(data.message); }
                } catch(e) {}
            }
        }

        function copyTxt(elementId, btnId) {
            let elem = document.getElementById(elementId);
            if(!elem) return;
            let urlText = elem.tagName === "SELECT" ? elem.value : elem.innerText;
            if(!urlText.includes("Menunggu") && !urlText.includes("Tidak Aktif")) {
                navigator.clipboard.writeText(urlText); let btn = document.getElementById(btnId);
                btn.innerText = "Done"; btn.style.background = "#1a1a1a"; btn.style.color = "#ffffff";
                setTimeout(() => { btn.innerText = "Copy"; btn.style.background = "#000000"; btn.style.color = "#ffffff"; }, 1500);
            }
        }

        async function fetchServerInfo() {
            try {
                const response = await fetch('/__info');
                if (!response.ok) return;
                const data = await response.json();
                if (data.uuid) document.getElementById('uuidInput').value = data.uuid;
                window.serverActivePaths = data.paths;
            } catch (e) {}
        }

        function buildConfig(protocol, type, evt) {
            document.querySelectorAll('.btn-blue').forEach(b => b.classList.remove('btn-active'));
            if(evt && evt.target) evt.target.classList.add('btn-active');
            const uuid = document.getElementById('uuidInput').value.trim();
            const hostSelect = document.getElementById('domainSelect');
            const host = hostSelect.value.trim();
            const bugHost = document.getElementById('bugInput').value.trim();
            const engine = document.getElementById('engineSelect').value;
            const area = document.getElementById('output-area');
            const label = document.getElementById('out-type');
            const txt = document.getElementById('configText');
            if (!host || host.includes("Menunggu")) {
                alert("Domain Tunnel belum siap / belum dipilih!");
                return;
            }
            const pathsMapping = window.serverActivePaths || { vless: '/vless-argo', vmess: '/vmess-argo', trojan: '/trojan-argo' };
            let basePath = pathsMapping[protocol] || '/' + protocol + '-argo';
            let remark = 'DDFATHU-' + protocol.toUpperCase() + '-' + type.toUpperCase() + '-' + engine.toUpperCase();
            let configResult = '';
            label.innerText = remark;
            function safeBtoa(str) {
                return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, function(match, p1) {
                    return String.fromCharCode('0x' + p1);
                }));
            }
            let netType = engine;
            if (type === 'sni') {
                if (protocol === 'vless') {
                    configResult = 'vless://' + uuid + '@' + bugHost + ':443?encryption=none&security=tls&sni=' + host + '&fp=randomized&type=' + netType + (netType === 'grpc' ? '&serviceName=grpc-service' : (netType === 'h2' ? '&path=/h2-path' : '&host=' + host + '&path=' + encodeURIComponent(basePath))) + '#' + encodeURIComponent(remark);
                } else if (protocol === 'vmess') {
                    let jsonVmess = { v: "2", ps: remark, add: bugHost, port: 443, id: uuid, aid: 0, scy: "auto", net: netType, type: netType === 'grpc' ? 'multi' : 'none', host: host, path: netType === 'grpc' ? 'grpc-service' : (netType === 'h2' ? '/h2-path' : basePath), tls: "tls", sni: host };
                    configResult = 'vmess://' + safeBtoa(JSON.stringify(jsonVmess));
                } else if (protocol === 'trojan') {
                    configResult = 'trojan://' + uuid + '@' + bugHost + ':443?security=tls&sni=' + host + '&type=' + netType + (netType === 'grpc' ? '&serviceName=grpc-service' : (netType === 'h2' ? '&path=/h2-path' : '&host=' + host + '&path=' + encodeURIComponent(basePath))) + '#' + encodeURIComponent(remark);
                }
            } else if (type === 'sni_reverse') {
                if (protocol === 'vless') {
                    configResult = 'vless://' + uuid + '@' + host + ':443?encryption=none&security=tls&sni=' + bugHost + '&fp=randomized&type=' + netType + (netType === 'grpc' ? '&serviceName=grpc-service' : (netType === 'h2' ? '&path=/h2-path' : '&host=' + bugHost + '&path=' + encodeURIComponent(basePath))) + '#' + encodeURIComponent(remark);
                } else if (protocol === 'vmess') {
                    let jsonVmess = { v: "2", ps: remark, add: host, port: 443, id: uuid, aid: 0, scy: "auto", net: netType, type: netType === 'grpc' ? 'multi' : 'none', host: bugHost, path: netType === 'grpc' ? 'grpc-service' : (netType === 'h2' ? '/h2-path' : basePath), tls: "tls", sni: bugHost };
                    configResult = 'vmess://' + safeBtoa(JSON.stringify(jsonVmess));
                } else if (protocol === 'trojan') {
                    configResult = 'trojan://' + uuid + '@' + host + ':443?security=tls&sni=' + bugHost + '&type=' + netType + (netType === 'grpc' ? '&serviceName=grpc-service' : (netType === 'h2' ? '&path=/h2-path' : '&host=' + bugHost + '&path=' + encodeURIComponent(basePath))) + '#' + encodeURIComponent(remark);
                }
            } else if (type === 'cdn') {
                let pathBug = '/' + bugHost + basePath;
                if (protocol === 'vless') {
                    configResult = 'vless://' + uuid + '@' + host + ':443?encryption=none&security=tls&sni=' + host + '&fp=randomized&type=' + netType + (netType === 'grpc' ? '&serviceName=grpc-service' : (netType === 'h2' ? '&path=/h2-path' : '&host=' + host + '&path=' + encodeURIComponent(pathBug))) + '#' + encodeURIComponent(remark);
                } else if (protocol === 'vmess') {
                    let jsonVmess = { v: "2", ps: remark, add: host, port: 443, id: uuid, aid: 0, scy: "auto", net: netType, type: netType === 'grpc' ? 'multi' : 'none', host: host, path: netType === 'grpc' ? 'grpc-service' : (netType === 'h2' ? '/h2-path' : pathBug), tls: "tls", sni: host };
                    configResult = 'vmess://' + safeBtoa(JSON.stringify(jsonVmess));
                } else if (protocol === 'trojan') {
                    configResult = 'trojan://' + uuid + '@' + host + ':443?security=tls&sni=' + host + '&type=' + netType + (netType === 'grpc' ? '&serviceName=grpc-service' : (netType === 'h2' ? '&path=/h2-path' : '&host=' + host + '&path=' + encodeURIComponent(pathBug))) + '#' + encodeURIComponent(remark);
                }
            }
            txt.innerText = configResult;
            area.style.display = 'block';
        }

        function copyOutConfig() {
            navigator.clipboard.writeText(document.getElementById('configText').innerText);
            alert('Config Berhasil Disalin!');
        }

        setInterval(updateStats, 300000);
        updateStats();
        fetchAccounts();
        fetchServerInfo();
    </script>
</body>
</html>`
