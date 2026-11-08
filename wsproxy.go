package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

const wsMagic = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func runWsProxy() {
	ln, err := net.Listen("tcp", ":"+itoa(WSTargetPort))
	if err != nil {
		log.Println("[WSPROXY] Listen error:", err)
		return
	}
	log.Printf("[WSPROXY] Dynamic Proxy Active on Port %d", WSTargetPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleWsClient(conn)
	}
}

func handleWsClient(client net.Conn) {
	defer func() { recover() }()

	cfg := loadWsProxyCfg()
	if tcp, ok := client.(*net.TCPConn); ok {
		tcp.SetNoDelay(true)
		tcp.SetKeepAlive(true)
	}

	buf := make([]byte, 32*1024)
	n, err := client.Read(buf)
	if err != nil || n == 0 {
		client.Close()
		return
	}
	raw := string(buf[:n])
	rawLower := strings.ToLower(raw)
	headers := parseHTTPHeaders(raw)

	isWs := strings.Contains(rawLower, "upgrade: websocket") || headers["upgrade"] == "websocket"

	wsKey := headers["sec-websocket-key"]
	if wsKey == "" {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		wsKey = base64.StdEncoding.EncodeToString(b)
	}
	h := sha1.New()
	h.Write([]byte(wsKey + wsMagic))
	acceptKey := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var response string
	if isWs {
		response = "HTTP/1.1 101 Switching Protocols\r\n" +
			"Upgrade: websocket\r\n" +
			"Connection: Upgrade\r\n" +
			"Sec-WebSocket-Accept: " + acceptKey + "\r\n"
		if proto := headers["sec-websocket-protocol"]; proto != "" {
			response += "Sec-WebSocket-Protocol: " + proto + "\r\n"
		}
		response += "\r\n"
	} else {
		response = "HTTP/1.1 101 Switching Protocols\r\n\r\n"
	}
	if _, err := client.Write([]byte(response)); err != nil {
		client.Close()
		return
	}

	sshConn, err := net.Dial("tcp", "127.0.0.1:"+itoa(cfg.SSHPort))
	if err != nil {
		client.Close()
		return
	}
	if tcp, ok := sshConn.(*net.TCPConn); ok {
		tcp.SetNoDelay(true)
		tcp.SetKeepAlive(true)
	}

	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			client.Close()
			sshConn.Close()
		})
	}

	// SSH -> client
	go func() {
		_, _ = io.Copy(client, sshConn)
		cleanup()
	}()

	// Client -> SSH dengan filter "SSH-" magic
	needle := []byte("SSH-")
	pending := []byte{}
	handshakeDone := false

	for {
		chunk := make([]byte, 32*1024)
		rn, err := client.Read(chunk)
		if err != nil {
			cleanup()
			return
		}
		if handshakeDone {
			if _, err := sshConn.Write(chunk[:rn]); err != nil {
				cleanup()
				return
			}
			continue
		}
		pending = append(pending, chunk[:rn]...)
		if idx := bytes.Index(pending, needle); idx >= 0 {
			handshakeDone = true
			if _, err := sshConn.Write(pending[idx:]); err != nil {
				cleanup()
				return
			}
			pending = nil
		} else if len(pending) > cfg.MaxBuffer {
			cleanup()
			return
		}
	}
}

func parseHTTPHeaders(raw string) map[string]string {
	h := map[string]string{}
	for i, line := range strings.Split(raw, "\r\n") {
		if i == 0 || line == "" {
			continue
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			k := strings.ToLower(strings.TrimSpace(line[:idx]))
			v := strings.TrimSpace(line[idx+1:])
			h[k] = v
		}
	}
	return h
}
