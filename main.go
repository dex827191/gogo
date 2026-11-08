package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.Ltime)
	log.Println("[BOOT] gogo Panel — Go Edition")

	initConfigDir()
	initPaths()

	go runMux()      // Port 8881 — TLS/WS traffic splitter
	go runWsProxy()  // Port 8880 — WebSocket -> SSH
	go statsLoop()   // Collector stats tiap 4 detik
	go runServer()   // Port 8081 — UI + HTTP API

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("[SHUTDOWN] Signal received, exiting...")
	time.Sleep(500 * time.Millisecond)
}
