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

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
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

	for {
		n, err := io.Copy(ws, os.Stdin)
		if err != nil {
			logger.Printf("error while copying %v", err)
		}
		if n == 0 {
			break
		}
	}
}
