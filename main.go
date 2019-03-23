// udptee writes from stdin to an udp server and then to stdout
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
	filterEmails := flag.Bool("filter-emails", false, "mask emails's handle with * before writing to udp server")
	flag.Parse()

	conn := MustUDPConnect(*host, *port)
	defer conn.Close()
	fudp := &Filter{delegate: conn, filterColors: *filterColors, filterEmails: *filterEmails}
	ws := io.MultiWriter(os.Stdout, fudp)

	// 32k buffer
	size := 32 * 1024
	buf := make([]byte, size)

	for {
		nr, err := os.Stdin.Read(buf)

		if nr > 0 {
			_, err := ws.Write(buf[0:nr])
			if err != nil {
				logger.Printf("error while writing: %v", err)
			}
		}

		if err != nil {
			// EOF means the piped program exited and closed its stdout
			if err == io.EOF {
				break
			} else {
				logger.Printf("error while reading: %v", err)
				os.Exit(1)
			}
		}
	}
}
