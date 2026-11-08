package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	cachedDiskUsage   = "38%"
	cachedSSHOnline   = "0"
	cachedUserDetails = "Semua user offline"
	statsMu           sync.Mutex
)

func runServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/__info", handleInfo)
	mux.HandleFunc("/api/logtunnel", handleLogTunnel)
	mux.HandleFunc("/api/lognamed", handleLogNamed)
	mux.HandleFunc("/api/setup-pass", handleSetupPass)
	mux.HandleFunc("/api/set-network", handleSetNetwork)
	mux.HandleFunc("/api/set-wsproxy", handleSetWsProxy)
	mux.HandleFunc("/api/set-system", handleSetSystem)
	mux.HandleFunc("/api/set-token", handleSetToken)
	mux.HandleFunc("/api/set-cfip", handleSetCfip)
	mux.HandleFunc("/api/add", handleAdd)
	mux.HandleFunc("/api/delete", handleDelete)
	mux.HandleFunc("/api/list", handleList)
	mux.HandleFunc("/api/login", handleLogin)
	mux.HandleFunc("/api/stats", handleStats)
	mux.HandleFunc("/ssh-ws", handleSSHWS)
	mux.HandleFunc("/"+SubPathEnv, handleSub)

	// Bootstrap: gen config + download xray/cloudflared
	go func() {
		regenerateConfig()
		downloadFilesAndRun()
	}()

	// Auto detect quick tunnel domain
	go func() {
		for {
			extractQuickDomain()
			time.Sleep(3 * time.Second)
		}
	}()

	log.Printf("[SERVER] UI & Gateway running on :%d", MainPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", MainPort), mux); err != nil {
		log.Fatal(err)
	}
}

// ============================================================
// DOWNLOAD & START XRAY + CLOUDFLARED
// ============================================================
func downloadFilesAndRun() {
	arch := runtime.GOARCH
	var webURL, botURL string
	if arch == "arm64" || arch == "arm" {
		webURL = "https://arm64.ssss.nyc.mn/web"
		botURL = "https://arm64.ssss.nyc.mn/bot"
	} else {
		webURL = "https://amd64.ssss.nyc.mn/web"
		botURL = "https://amd64.ssss.nyc.mn/bot"
	}

	if err := downloadFile(webPath, webURL); err != nil {
		log.Printf("[DL] web error: %v", err)
		return
	}
	if err := downloadFile(botPath, botURL); err != nil {
		log.Printf("[DL] bot error: %v", err)
		return
	}
	_ = os.Chmod(webPath, 0775)
	_ = os.Chmod(botPath, 0775)

	startXray()

	// Tunnel: check saved token
	if b, err := os.ReadFile(ZTSingleTokenFile); err == nil {
		if tok := strings.TrimSpace(string(b)); tok != "" {
			restartSingleTunnel(tok)
			return
		}
	}
	restartSingleTunnel("")
}

func downloadFile(fileName, fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func startXray() {
	_ = exec.Command("pkill", "-9", "-f", webPath).Run()
	time.Sleep(300 * time.Millisecond)

	cmd := exec.Command(webPath, "-c", ConfigJSON)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		log.Printf("[XRAY] start error: %v", err)
	}
}

func restartXray() {
	startXray()
}

// ============================================================
// HANDLERS
// ============================================================
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

func handleSub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	b, err := os.ReadFile(SubFile)
	if err != nil {
		_, _ = w.Write([]byte("Loading sub..."))
		return
	}
	_, _ = w.Write(b)
}

func handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vless := readPathsFile("pathvless.txt", "/vless-argo")[0]
	vmess := readPathsFile("pathvmess.txt", "/vmess-argo")[0]
	trojan := readPathsFile("pathtrojan.txt", "/trojan-argo")[0]
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"uuid":   UUID,
		"cfip":   getActiveCfip(),
		"domain": currentActiveDomain,
		"paths":  map[string]string{"vless": vless, "vmess": vmess, "trojan": trojan},
	})
}

func handleLogTunnel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	b, err := os.ReadFile(LogPath)
	if err != nil {
		_, _ = w.Write([]byte("Log belum siap."))
		return
	}
	_, _ = w.Write(b)
}

func handleLogNamed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	b, err := os.ReadFile(ZTLogPath)
	if err != nil {
		_, _ = w.Write([]byte("Log Zero Trust belum terbuat."))
		return
	}
	_, _ = w.Write(b)
}

func handleSetupPass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	newPass := strings.TrimSpace(q.Get("pass"))
	oldPass := strings.TrimSpace(q.Get("old_pass"))
	cur := getAdminPassword()

	if cur != "" && oldPass != cur {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Password Admin Lama Salah!"})
		return
	}
	if len(newPass) < 4 {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Password minimal 4 karakter!"})
		return
	}
	_ = os.WriteFile(AdminPassFile, []byte(newPass), 0600)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Password Admin Berhasil Disimpan!"})
}

func handleSetNetwork(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if !verifyAdminPassword(q.Get("pass")) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Akses Ditolak!"})
		return
	}
	n := NetSettings{
		DnsType:   getOr(q.Get("dns_type"), "udp"),
		CustomDNS: getOr(q.Get("custom_dns"), "8.8.8.8"),
		Engine:    getOr(q.Get("engine"), "ws"),
	}
	saveNetSettings(n)
	regenerateConfig()
	restartXray()
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("V2Ray Config Disimpan! DNS:[%s] %s Engine:%s", n.DnsType, n.CustomDNS, n.Engine),
	})
}

func handleSetWsProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if !verifyAdminPassword(q.Get("pass")) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Akses Ditolak!"})
		return
	}
	c := WsProxyConfig{
		SSHPort:   atoi(getOr(q.Get("ssh_port"), "22")),
		KeepAlive: atoi(getOr(q.Get("keep_alive"), "15000")),
		MaxBuffer: atoi(getOr(q.Get("max_buffer"), "32768")),
	}
	saveWsProxyCfg(c)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Setting SSH Disimpan!"})
}

func handleSetSystem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if !verifyAdminPassword(q.Get("pass")) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Akses Ditolak!"})
		return
	}
	c := SystemCfg{
		Banner:          q.Get("banner"),
		EnableBBR:       getOr(q.Get("enable_bbr"), "true"),
		UdpgwPort:       getOr(q.Get("udpgw_port"), "7300"),
		UdpgwMaxClients: getOr(q.Get("udpgw_max_clients"), "1000"),
	}
	saveSystemCfg(c)

	if c.Banner != "" {
		_ = os.WriteFile("/etc/dropbear_banner", []byte(c.Banner), 0644)
	} else {
		def := "==================================================\n" +
			"                 SSH SERVER PANEL                 \n" +
			"==================================================\n" +
			" powered by : d e d e f a t h u\n" +
			"==================================================\n"
		_ = os.WriteFile("/etc/dropbear_banner", []byte(def), 0644)
	}

	ws := loadWsProxyCfg()
	_ = exec.Command("pkill", "-9", "dropbear").Run()
	time.Sleep(500 * time.Millisecond)
	cmd := exec.Command("/usr/sbin/dropbear",
		"-p", fmt.Sprintf("127.0.0.1:%d", ws.SSHPort),
		"-b", "/etc/dropbear_banner",
		"-W", "1048576", "-K", "15", "-I", "300")
	_ = cmd.Start()

	_ = exec.Command("pkill", "-9", "badvpn-udpgw").Run()
	time.Sleep(300 * time.Millisecond)
	if _, err := os.Stat("/usr/local/bin/badvpn-udpgw"); err == nil {
		cmd2 := exec.Command("/usr/local/bin/badvpn-udpgw",
			"--listen-addr", "0.0.0.0:"+c.UdpgwPort,
			"--max-clients", c.UdpgwMaxClients,
			"--max-connections-for-client", "50")
		_ = cmd2.Start()
	}

	cong := "cubic"
	if c.EnableBBR == "true" {
		cong = "bbr"
	}
	_ = exec.Command("sysctl", "-w", "net.ipv4.tcp_congestion_control="+cong).Run()

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "System Config Berhasil Disimpan!"})
}

func handleSetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if !verifyAdminPassword(q.Get("pass")) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Akses Ditolak!"})
		return
	}
	go restartSingleTunnel(q.Get("token"))
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Perintah restart tunnel terkirim!"})
}

func handleSetCfip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if !verifyAdminPassword(q.Get("pass")) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Akses Ditolak!"})
		return
	}
	ip := strings.TrimSpace(q.Get("ip"))
	if ip != "" {
		_ = os.WriteFile(CfipFile, []byte(ip), 0644)
		if currentActiveDomain != "" {
			generateLinks(currentActiveDomain)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "CFIP diperbarui ke " + ip})
	} else {
		_ = os.Remove(CfipFile)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "CFIP direset ke default."})
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	ip := r.Header.Get("CF-Connecting-IP")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	ua := r.Header.Get("User-Agent")
	ok, msg := addSshUser(q.Get("user"), q.Get("pass"), ip, ua, getCurrentHosts())
	if !ok {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": msg})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": msg})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if !verifyAdminPassword(q.Get("token")) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Akses Ditolak!"})
		return
	}
	ok, msg := deleteSshUser(q.Get("user"))
	if !ok {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": msg})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": msg})
}

func handleList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users := listSshUsers()
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"total":  len(users),
		"users":  users,
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if getAdminPassword() == "" {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "not_configured"})
		return
	}
	p := r.URL.Query().Get("pass")
	if verifyAdminPassword(p) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "token": p})
	} else {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Password Salah!"})
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statsMu.Lock()
	disk, sshOn, details := cachedDiskUsage, cachedSSHOnline, cachedUserDetails
	statsMu.Unlock()

	railwayURL := os.Getenv("SNI")
	if railwayURL == "" {
		railwayURL = "Tidak Aktif"
	}
	if rt := os.Getenv("RAILWAY_TCP_PROXY_DOMAIN"); rt != "" {
		railwayURL = rt + ":" + os.Getenv("RAILWAY_TCP_PROXY_PORT")
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"cpu_model":         getCPUModel(),
		"ram_total":         fmt.Sprintf("%.2f GB", float64(memTotal())/1073741824.0),
		"ram_used":          fmt.Sprintf("%.2f GB", float64(memTotal()-memAvailable())/1073741824.0),
		"disk_usage":        disk,
		"uptime":            fmt.Sprintf("%.2f Hours", uptimeHours()),
		"ssh_online":        sshOn,
		"user_list_details": details,
		"active_cfip":       getActiveCfip(),
		"quick_url":         currentActiveDomain,
		"zt_domains":        getDomainsByPort([]string{"8880", "8881"}),
		"zt_vmess_domains":  getDomainsByPort([]string{"8001"}),
		"pass_configured":   getAdminPassword() != "",
		"railway_url":       railwayURL,
		"status":            "ONLINE",
		"dns_type":          loadNetSettings().DnsType,
		"custom_dns":        loadNetSettings().CustomDNS,
		"engine_mode":       loadNetSettings().Engine,
		"ws_proxy_cfg":      loadWsProxyCfg(),
		"sys_settings":      loadSystemCfg(),
	})
}

// SSH over WS via HTTP hijack (fallback untuk /ssh-ws di port 8081)
func handleSSHWS(w http.ResponseWriter, r *http.Request) {
	ws := loadWsProxyCfg()
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijack not supported", 500)
		return
	}
	clientConn, bufrw, err := hj.Hijack()
	if err != nil {
		return
	}
	defer clientConn.Close()

	target, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", ws.SSHPort))
	if err != nil {
		return
	}
	defer target.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.URL.RequestURI(), r.Proto))
	for k, vv := range r.Header {
		for _, v := range vv {
			sb.WriteString(k + ": " + v + "\r\n")
		}
	}
	sb.WriteString("\r\n")
	_, _ = target.Write([]byte(sb.String()))

	go func() { _, _ = io.Copy(target, bufrw) }()
	_, _ = io.Copy(clientConn, target)
}

// ============================================================
// HELPERS — CPU / RAM / DISK / UPTIME
// ============================================================
func getCPUModel() string {
	b, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return runtime.GOARCH
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "model name") {
			if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return runtime.GOARCH
}

func memTotal() uint64 {
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return uint64(atoi(fields[1])) * 1024
			}
		}
	}
	return 0
}

func memAvailable() uint64 {
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				return uint64(atoi(fields[1])) * 1024
			}
		}
	}
	return 0
}

func uptimeHours() float64 {
	b, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(b))
	if len(fields) < 1 {
		return 0
	}
	sec := 0.0
	fmt.Sscanf(fields[0], "%f", &sec)
	return sec / 3600.0
}

// ============================================================
// STATS LOOP
// ============================================================
func statsLoop() {
	for {
		if out, err := exec.Command("df", "-h", "/").Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) > 1 {
				fields := strings.Fields(lines[1])
				if len(fields) >= 5 {
					statsMu.Lock()
					cachedDiskUsage = fields[4]
					statsMu.Unlock()
				}
			}
		}

		out, err := exec.Command("sh", "-c",
			"netstat -anp 2>/dev/null | grep dropbear | grep ESTABLISHED | awk '{print $5}' | cut -d: -f1 | sort -u").Output()
		if err == nil && strings.TrimSpace(string(out)) != "" {
			ips := strings.Fields(string(out))
			statsMu.Lock()
			cachedSSHOnline = itoa(len(ips))
			if len(ips) > 0 {
				cachedUserDetails = "IP Active: " + ips[0]
			}
			statsMu.Unlock()
		} else {
			statsMu.Lock()
			cachedSSHOnline = "0"
			cachedUserDetails = "Semua user offline"
			statsMu.Unlock()
		}
		time.Sleep(4 * time.Second)
	}
}

// ============================================================
// XRAY CONFIG GENERATOR
// ============================================================
func regenerateConfig() {
	n := loadNetSettings()
	vlessPaths := readPathsFile("pathvless.txt", "/vless-argo")
	vmessPaths := readPathsFile("pathvmess.txt", "/vmess-argo")
	trojanPaths := readPathsFile("pathtrojan.txt", "/trojan-argo")

	fallbacks := []map[string]interface{}{}
	inbounds := []map[string]interface{}{}
	port := 3100

	for _, p := range vlessPaths {
		fallbacks = append(fallbacks, map[string]interface{}{"path": p, "dest": port})
		inbounds = append(inbounds, buildInbound(port, "vless", p, n.Engine))
		port++
	}
	for _, p := range vmessPaths {
		fallbacks = append(fallbacks, map[string]interface{}{"path": p, "dest": port})
		inbounds = append(inbounds, buildInbound(port, "vmess", p, n.Engine))
		port++
	}
	for _, p := range trojanPaths {
		fallbacks = append(fallbacks, map[string]interface{}{"path": p, "dest": port})
		inbounds = append(inbounds, buildInbound(port, "trojan", p, n.Engine))
		port++
	}

	mainInbound := map[string]interface{}{
		"port":     ArgoPort,
		"protocol": "vless",
		"settings": map[string]interface{}{
			"clients":    []map[string]string{{"id": UUID}},
			"decryption": "none",
			"fallbacks":  fallbacks,
		},
		"streamSettings": map[string]string{"network": "tcp", "security": "none"},
		"sniffing":       map[string]interface{}{"enabled": true, "destOverride": []string{"http", "tls", "quic"}},
	}
	inbounds = append([]map[string]interface{}{mainInbound}, inbounds...)

	dns := []string{n.CustomDNS, "1.1.1.1"}
	cfg := map[string]interface{}{
		"log":       map[string]string{"access": "/dev/null", "error": "/dev/null", "loglevel": "none"},
		"inbounds":  inbounds,
		"dns":       map[string]interface{}{"servers": dns},
		"outbounds": []map[string]string{{"protocol": "freedom", "tag": "direct"}},
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(ConfigJSON, b, 0644)
}

func buildInbound(port int, proto, path, engine string) map[string]interface{} {
	stream := map[string]interface{}{"network": engine, "security": "none"}
	if engine == "ws" {
		stream["wsSettings"] = map[string]string{"path": path}
	}
	var settings map[string]interface{}
	switch proto {
	case "vless":
		settings = map[string]interface{}{"clients": []map[string]interface{}{{"id": UUID, "level": 0}}, "decryption": "none"}
	case "vmess":
		settings = map[string]interface{}{"clients": []map[string]interface{}{{"id": UUID, "alterId": 0}}}
	case "trojan":
		settings = map[string]interface{}{"clients": []map[string]interface{}{{"password": UUID}}}
	}
	return map[string]interface{}{
		"port":           port,
		"listen":         "127.0.0.1",
		"protocol":       proto,
		"settings":       settings,
		"streamSettings": stream,
		"sniffing":       map[string]interface{}{"enabled": true, "destOverride": []string{"http", "tls", "quic"}},
	}
}

// ============================================================
// SUB GENERATOR
// ============================================================
func generateLinks(domain string) {
	cfip := getActiveCfip()
	vlessPath := readPathsFile("pathvless.txt", "/vless-argo")[0]
	vmessPath := readPathsFile("pathvmess.txt", "/vmess-argo")[0]
	trojanPath := readPathsFile("pathtrojan.txt", "/trojan-argo")[0]
	node := Name + "-Server"

	vmessJSON := fmt.Sprintf(`{"v":"2","ps":"%s","add":"%s","port":"%s","id":"%s","aid":"0","scy":"auto","net":"ws","type":"none","host":"%s","path":"%s?ed=2560","tls":"tls","sni":"%s","alpn":"","fp":"firefox"}`,
		node, cfip, Cfport, UUID, domain, vmessPath, domain)

	sub := fmt.Sprintf(
		"vless://%s@%s:%s?encryption=none&security=tls&sni=%s&fp=firefox&type=ws&host=%s&path=%s#%s\n\n"+
			"vmess://%s\n\n"+
			"trojan://%s@%s:%s?security=tls&sni=%s&fp=firefox&type=ws&host=%s&path=%s#%s",
		UUID, cfip, Cfport, domain, domain, urlEncodePath(vlessPath+"?ed=2560"), node,
		base64.StdEncoding.EncodeToString([]byte(vmessJSON)),
		UUID, cfip, Cfport, domain, domain, urlEncodePath(trojanPath+"?ed=2560"), node,
	)
	_ = os.WriteFile(SubFile, []byte(base64.StdEncoding.EncodeToString([]byte(sub))), 0644)
}

func urlEncodePath(s string) string {
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == '~' || r == '/' || r == '?' || r == '=' || r == '&' {
			b.WriteRune(r)
		} else {
			b.WriteString(fmt.Sprintf("%%%02X", r))
		}
	}
	return b.String()
}
