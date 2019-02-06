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

func runOnce(writer io.Writer, reader io.Reader) bool {
	n, err := io.Copy(writer, reader)
	if err != nil {
		logger.Printf("error while copying %v", err)
	}
	if n == 0 {
		return false
	}
	return true
}

func connectAndRun(host string, port int, filterColors bool, filterEmails bool, writer io.Writer, reader io.Reader, times int) {
	conn := MustUDPConnect(host, port)
	defer conn.Close()
	fudp := &Filter{delegate: conn, filterColors: filterColors, filterEmails: filterEmails}
	ws := io.MultiWriter(writer, fudp)

	if times == 0 {
		for {
			if !runOnce(ws, reader) {
				break
			}
		}
	} else {
		for i := 0; i < times; i++ {
			if !runOnce(ws, reader) {
				break
			}
		}
	}
}

func main() {
	host := flag.String("host", defaultHost, "udp host")
	port := flag.Int("port", defaultPort, "udp port")
	filterColors := flag.Bool("filter-colors", false, "remove colors from log before writing to udp server")
	filterEmails := flag.Bool("filter-emails", false, "mask emails with * before writing to udp server")
	flag.Parse()

	connectAndRun(*host, *port, *filterColors, *filterEmails, os.Stdout, os.Stdin, 0)
}
