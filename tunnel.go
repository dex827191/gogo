package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var currentActiveDomain string

func restartSingleTunnel(token string) {
	_ = exec.Command("pkill", "-9", "-f", "cloudflared").Run()

	bin, _ := exec.LookPath("cloudflared")
	if bin == "" {
		bin = "cloudflared"
	}

	if strings.TrimSpace(token) != "" {
		_ = os.WriteFile(ZTSingleTokenFile, []byte(strings.TrimSpace(token)), 0644)
		cmd := exec.Command(bin, "tunnel", "run",
			"--protocol", "http2", "--no-tls-verify",
			"--token", strings.TrimSpace(token))
		if f, err := os.OpenFile(ZTLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			cmd.Stdout = f
			cmd.Stderr = f
		}
		_ = cmd.Start()
	} else {
		_ = os.Remove(ZTSingleTokenFile)
		_ = os.WriteFile(ZTLogPath, []byte("Token removed."), 0644)
		cmd := exec.Command(bin, "tunnel", "--edge-ip-version", "auto",
			"--no-autoupdate", "--protocol", "http2",
			"--logfile", LogPath, "--loglevel", "info",
			"--url", "http://localhost:8001")
		_ = cmd.Start()
	}
}

type DomainEntry struct {
	Domain string `json:"domain"`
	Port   string `json:"port"`
}

func getDomainsByPort(targetPorts []string) []DomainEntry {
	domains := []DomainEntry{}
	b, err := os.ReadFile(ZTLogPath)
	if err != nil {
		return domains
	}
	content := string(b)
	portPattern := strings.Join(targetPorts, "|")

	re1 := regexp.MustCompile(`(?s)"hostname"\s*:\s*"([^"]+)".*?localhost:(` + portPattern + `)`)
	for _, m := range re1.FindAllStringSubmatch(content, -1) {
		if !containsDomain(domains, m[1]) {
			domains = append(domains, DomainEntry{Domain: m[1], Port: m[2]})
		}
	}
	if len(domains) == 0 {
		re2 := regexp.MustCompile(`(?s)localhost:(` + portPattern + `).*?"hostname"\s*:\s*"([^"]+)"`)
		for _, m := range re2.FindAllStringSubmatch(content, -1) {
			if !containsDomain(domains, m[2]) {
				domains = append(domains, DomainEntry{Domain: m[2], Port: m[1]})
			}
		}
	}
	return domains
}

func containsDomain(list []DomainEntry, d string) bool {
	for _, x := range list {
		if x.Domain == d {
			return true
		}
	}
	return false
}

func getCurrentHosts() string {
	out := ""
	sshDoms := getDomainsByPort([]string{"8880", "8881"})
	namedURL := ""
	if len(sshDoms) > 0 {
		namedURL = sshDoms[0].Domain
	} else if d := os.Getenv("D"); d != "" {
		namedURL = d
	}
	if namedURL != "" && !strings.Contains(namedURL, "Menghubungkan") {
		out += strings.TrimPrefix(strings.TrimPrefix(namedURL, "https://"), "http://") + " (SSH WS)"
	}
	if rt := os.Getenv("RAILWAY_TCP_PROXY_DOMAIN"); rt != "" {
		if port := os.Getenv("RAILWAY_TCP_PROXY_PORT"); port != "" {
			if out != "" {
				out += " dan ini (SSH SNI) " + rt + ":" + port
			} else {
				out = rt + ":" + port
			}
		}
	} else if sni := os.Getenv("SNI"); sni != "" {
		sni = strings.TrimPrefix(strings.TrimPrefix(sni, "https://"), "http://")
		if out != "" {
			out += " dan " + sni
		} else {
			out = sni
		}
	}
	if out == "" {
		out = "Menunggu Quick Tunnel..."
	}
	return out
}

func extractQuickDomain() {
	b, err := os.ReadFile(LogPath)
	if err != nil {
		return
	}
	re := regexp.MustCompile(`https://([a-zA-Z0-9\-]+\.trycloudflare\.com)`)
	if m := re.FindStringSubmatch(string(b)); len(m) == 2 {
		if currentActiveDomain != m[1] {
			currentActiveDomain = m[1]
			generateLinks(currentActiveDomain)
		}
	}
}
