package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// scanner := bufio.NewScanner(os.Stdin)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	f, err := os.OpenFile("log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	flog := &Filter{delegate: f, filterColors: true, filterEmails: true}
	fudp := &Filter{delegate: conn, filterColors: true, filterEmails: true}
	ws := io.MultiWriter(os.Stdout, flog, fudp)

	if _, err := io.Copy(ws, os.Stdin); err != nil {
		log.Fatalf("error while copying %v", err)
	}
}
