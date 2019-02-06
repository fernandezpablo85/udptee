package main

import (
	"testing"
)

func TestUDPConnectOK(t *testing.T) {
	_, err := udpConnect("localhost", 1234)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestUDPConnectNoSuchHost(t *testing.T) {
	_, err := udpConnect("no_such_host", 1234)
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestUDPConnectNegativePort(t *testing.T) {
	_, err := udpConnect("localhost", -1)
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestUDPConnectHighPort(t *testing.T) {
	_, err := udpConnect("localhost", 2<<15)
	if err == nil {
		t.Errorf("expected error but got none %v", err)
	}
}
