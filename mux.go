package main

import (
	"io"
	"log"
	"net"
)

func runMux() {
	ln, err := net.Listen("tcp", ":"+itoa(PublicPort))
	if err != nil {
		log.Println("[MUX] Listen error:", err)
		return
	}
	log.Printf("[MUX] Pure Traffic Splitter Active on Port %d", PublicPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleMuxClient(conn)
	}
}

func handleMuxClient(client net.Conn) {
	defer func() { recover() }()

	if tcp, ok := client.(*net.TCPConn); ok {
		tcp.SetNoDelay(true)
		tcp.SetKeepAlive(true)
	}

	buf := make([]byte, 64*1024)
	n, err := client.Read(buf)
	if err != nil || n == 0 {
		client.Close()
		return
	}

	targetPort := WSTargetPort
	if buf[0] == 0x16 { // TLS handshake byte
		targetPort = SSLTargetPort
	}

	backend, err := net.Dial("tcp", "127.0.0.1:"+itoa(targetPort))
	if err != nil {
		client.Close()
		return
	}
	if btcp, ok := backend.(*net.TCPConn); ok {
		btcp.SetNoDelay(true)
		btcp.SetKeepAlive(true)
	}

	if _, err := backend.Write(buf[:n]); err != nil {
		client.Close()
		backend.Close()
		return
	}

	done := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(backend, client); done <- struct{}{} }()
	go func() { _, _ = io.Copy(client, backend); done <- struct{}{} }()
	<-done
	client.Close()
	backend.Close()
}
