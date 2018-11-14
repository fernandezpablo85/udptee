package main

import (
	"flag"
	"io"
	"log"
	"os"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 4321
)

func main() {
	host := flag.String("host", defaultHost, "udp host")
	port := flag.Int("port", defaultPort, "udp port")
	filterColors := flag.Bool("filter-colors", false, "remove colors from log before writing to udp server")
	filterEmails := flag.Bool("filter-emails", false, "mask emails with * before writing to udp server")
	flag.Parse()

	conn := MustUDPConnect(*host, *port)
	defer conn.Close()
	fudp := &Filter{delegate: conn, filterColors: *filterColors, filterEmails: *filterEmails}
	ws := io.MultiWriter(os.Stdout, fudp)

	if _, err := io.Copy(ws, os.Stdin); err != nil {
		log.Fatalf("error while copying %v", err)
	}
}
