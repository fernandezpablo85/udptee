package main

import (
	"fmt"
	"log"
	"net"
)

// UDPConnect connects to the given host and port
func UDPConnect(host string, port int) (net.Conn, error) {
	connstr := fmt.Sprintf("%s:%d", host, port)
	return net.Dial("udp", connstr)
}

// MustUDPConnect connects to the given host and port or crashes the program
func MustUDPConnect(host string, port int) net.Conn {
	conn, err := UDPConnect(host, port)
	if err != nil {
		log.Fatalf("error while connecting to %s:%d: %v", host, port, err)
	}
	return conn
}
