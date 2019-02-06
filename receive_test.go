package main

import (
	"bytes"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TestReceiveUDPMessage(t *testing.T) {
	pc, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Errorf("expected to listen udp, got %v", err)
	}
	defer pc.Close()
	port, err := strconv.ParseUint(strings.Split(pc.LocalAddr().String(), ":")[1], 10, 16)
	if err != nil {
		t.Errorf("invalid port assigned, %v", err)
	}
	writer := new(bytes.Buffer)
	msg := []byte("hello\n")
	reader := bytes.NewReader(msg)
	connectAndRun("127.0.0.1", int(port), true, true, writer, reader, 1)
	if bytes.Compare(writer.Bytes(), msg) != 0 {
		t.Errorf("expected writer to have received message, got %v", msg)
	}
	buf := make([]byte, len(msg))
	n, _, err := pc.ReadFrom(buf)
	if err != nil {
		t.Errorf("failed to read udp message: %v", err)
	}
	if n != len(msg) {
		t.Errorf("expected to read %d bytes, got %d", n, len(msg))
	}
	if bytes.Compare(writer.Bytes(), buf) != 0 {
		t.Errorf("expected server to have received message, got %v", buf)
	}
}
