package lookup

import (
	"net/http"
	"strings"
	"testing"
)

func TestIPAddrFromRemoteAddr(t *testing.T) {
	ip := "127.0.0.1:8080"
	ipLookups := []string{"RemoteAddr"}
	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	if err != nil {
		t.Errorf("Unable to create new HTTP request. Error: %v", err)
	}

	request.RemoteAddr = ip

	if a := RemoteIP(ipLookups, request); a != "127.0.0.1" {
		t.Errorf("ipAddrFromRemoteAddr did not chop the port number correctly. Got %v", a)
	}
}

func TestRemoteIPDefault(t *testing.T) {
	ipLookups := []string{"RemoteAddr", "X-Real-IP"}
	ipv6 := "2601:7:1c82:4097:59a0:a80b:2841:b8c8"

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	if err != nil {
		t.Errorf("Unable to create new HTTP request. Error: %v", err)
	}

	request.Header.Set("X-Real-IP", ipv6)

	ip := RemoteIP(ipLookups, request)
	if ip != request.RemoteAddr {
		t.Errorf("Did not get the right IP. IP: %v", ip)
	}
	if ip == ipv6 {
		t.Errorf("X-Real-IP should have been skipped. IP: %v", ip)
	}
}

func TestRemoteIPForwardedFor(t *testing.T) {
	ipLookups := []string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"}
	ipv6 := "2601:7:1c82:4097:59a0:a80b:2841:b8c8"

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	if err != nil {
		t.Errorf("Unable to create new HTTP request. Error: %v", err)
	}

	request.Header.Set("X-Forwarded-For", "54.223.11.104")
	request.Header.Set("X-Real-IP", ipv6)

	ip := RemoteIP(ipLookups, request)
	if ip != "54.223.11.104" {
		t.Errorf("Did not get the right IP. IP: %v", ip)
	}
	if ip == ipv6 {
		t.Errorf("X-Real-IP should have been skipped. IP: %v", ip)
	}
}

func TestRemoteIPRealIP(t *testing.T) {
	ipLookups := []string{"X-Real-IP", "X-Forwarded-For", "RemoteAddr"}
	ipv6 := "2601:7:1c82:4097:59a0:a80b:2841:b8c8"

	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	if err != nil {
		t.Errorf("Unable to create new HTTP request. Error: %v", err)
	}

	request.Header.Set("X-Forwarded-For", "54.223.11.104")
	request.Header.Set("X-Real-IP", ipv6)

	ip := RemoteIP(ipLookups, request)
	if ip != ipv6 {
		t.Errorf("Did not get the right IP. IP: %v", ip)
	}
	if ip == "54.223.11.104" {
		t.Errorf("X-Forwarded-For should have been skipped. IP: %v", ip)
	}
}

func TestRemoteIPEmpty(t *testing.T) {
	ipLookups := []string{"X-Real-IP", "X-Forwarded-For"}
	request, err := http.NewRequest("GET", "/", strings.NewReader("Hello, world!"))
	if err != nil {
		t.Errorf("Unable to create new HTTP request. Error: %v", err)
	}
	ip := RemoteIP(ipLookups, request)
	if ip != "" {
		t.Errorf("Expected empty response")
	}
}
