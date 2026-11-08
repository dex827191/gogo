package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Ports & constants
const (
	PublicPort    = 8881
	SSLTargetPort = 2443
	WSTargetPort  = 8880
	MainPort      = 8081
	ArgoPort      = 8001

	FilePath          = ".tmp"
	LogPath           = ".tmp/boot.log"
	ZTLogPath         = "/tmp/named_tunnel.log"
	ZTSingleTokenFile = "/tmp/zt_single_token.txt"
	CfipFile          = "/tmp/cfip.txt"
	NetSettingFile    = "/tmp/net_settings.json"
	WsProxyCfgFile    = "/tmp/ws_proxy_config.json"
	SystemCfgFile     = "/tmp/system_config.json"
	AdminPassFile     = "/tmp/admin_pass.txt"
	StatsPath         = "/tmp/server_stats.json"
	DBPath            = "/tmp/ssh_details.json"
	ConfigJSON        = ".tmp/config.json"
	SubFile           = ".tmp/sub.txt"

	UUID = "1f37ac4f-fdd0-49df-9406-1eda70a1d512"
)

var (
	Name       = getEnv("NAME", "ddfathu")
	Cfport     = getEnv("CFPORT", "443")
	SubPathEnv = getEnv("SUB_PATH", "sub")

	webPath string
	botPath string
)

func initPaths() {
	rand.Seed(time.Now().UnixNano())
	webPath = filepath.Join(FilePath, randomName())
	botPath = filepath.Join(FilePath, randomName())
}

func randomName() string {
	const chars = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 6)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func initConfigDir() {
	_ = os.MkdirAll(FilePath, 0755)
}

// ---------- Network Settings ----------
type NetSettings struct {
	DnsType   string `json:"dns_type"`
	CustomDNS string `json:"custom_dns"`
	Engine    string `json:"engine"`
}

func loadNetSettings() NetSettings {
	n := NetSettings{DnsType: "udp", CustomDNS: "8.8.8.8", Engine: "ws"}
	if b, err := os.ReadFile(NetSettingFile); err == nil {
		_ = json.Unmarshal(b, &n)
	}
	return n
}

func saveNetSettings(n NetSettings) {
	b, _ := json.MarshalIndent(n, "", "  ")
	_ = os.WriteFile(NetSettingFile, b, 0644)
}

// ---------- WS Proxy Config ----------
type WsProxyConfig struct {
	SSHPort   int `json:"sshPort"`
	KeepAlive int `json:"keepAlive"`
	MaxBuffer int `json:"maxBuffer"`
}

func loadWsProxyCfg() WsProxyConfig {
	c := WsProxyConfig{SSHPort: 22, KeepAlive: 15000, MaxBuffer: 32768}
	if b, err := os.ReadFile(WsProxyCfgFile); err == nil {
		_ = json.Unmarshal(b, &c)
	}
	return c
}

func saveWsProxyCfg(c WsProxyConfig) {
	b, _ := json.MarshalIndent(c, "", "  ")
	_ = os.WriteFile(WsProxyCfgFile, b, 0644)
}

// ---------- System Settings ----------
type SystemCfg struct {
	Banner          string `json:"banner"`
	EnableBBR       string `json:"enable_bbr"`
	UdpgwPort       string `json:"udpgw_port"`
	UdpgwMaxClients string `json:"udpgw_max_clients"`
}

func loadSystemCfg() SystemCfg {
	c := SystemCfg{EnableBBR: "true", UdpgwPort: "7300", UdpgwMaxClients: "1000"}
	if b, err := os.ReadFile(SystemCfgFile); err == nil {
		_ = json.Unmarshal(b, &c)
	}
	return c
}

func saveSystemCfg(c SystemCfg) {
	b, _ := json.MarshalIndent(c, "", "  ")
	_ = os.WriteFile(SystemCfgFile, b, 0644)
}

// ---------- CFIP ----------
func getActiveCfip() string {
	if b, err := os.ReadFile(CfipFile); err == nil {
		if s := string(b); s != "" {
			return s
		}
	}
	return getEnv("CFIP", "104.17.3.81")
}

// ---------- Admin Password ----------
func getAdminPassword() string {
	if b, err := os.ReadFile(AdminPassFile); err == nil {
		return string(b)
	}
	return os.Getenv("ADMIN_PASSWORD")
}

func verifyAdminPassword(p string) bool {
	cur := getAdminPassword()
	return cur != "" && cur == p
}

// ---------- SSH DB ----------
type UserMeta struct {
	Password  string `json:"password"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

func loadDB() map[string]UserMeta {
	db := map[string]UserMeta{}
	if b, err := os.ReadFile(DBPath); err == nil {
		_ = json.Unmarshal(b, &db)
	}
	return db
}

func saveDB(db map[string]UserMeta) {
	b, _ := json.MarshalIndent(db, "", "  ")
	_ = os.WriteFile(DBPath, b, 0644)
}

// ---------- Paths file helper ----------
func readPathsFile(name, def string) []string {
	if b, err := os.ReadFile(name); err == nil {
		out := []string{}
		for _, line := range splitLines(string(b)) {
			if len(line) > 0 && line[0] == '/' {
				out = append(out, line)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return []string{def}
}

func splitLines(s string) []string {
	var res []string
	cur := ""
	for _, r := range s {
		if r == '\n' {
			res = append(res, cur)
			cur = ""
		} else if r != '\r' {
			cur += string(r)
		}
	}
	if cur != "" {
		res = append(res, cur)
	}
	return res
}

// ---------- helpers ----------
func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return n
		}
		n = n*10 + int(r-'0')
	}
	return n
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var digits []byte
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	if neg {
		return "-" + string(digits)
	}
	return string(digits)
}

func getOr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
