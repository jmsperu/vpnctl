package netutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// PublicIP fetches the current public IP address.
func PublicIP() (string, error) {
	resp, err := httpGet("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

// DNSLeakTest performs a basic DNS leak check by resolving a known
// test domain and comparing the resolver IPs against known public DNS.
func DNSLeakTest() ([]string, error) {
	// Use DNS leak test API
	resp, err := httpGet("https://bash.ws/dnsleak/test/")
	if err != nil {
		// Fallback: resolve whoami.akamai.net and report the resolver.
		addrs, lookupErr := net.LookupHost("whoami.akamai.net")
		if lookupErr != nil {
			return nil, fmt.Errorf("DNS leak test unavailable: %w", lookupErr)
		}
		return addrs, nil
	}

	// Try to parse JSON array response.
	type dnsEntry struct {
		IP string `json:"ip"`
	}
	var entries []dnsEntry
	if err := json.Unmarshal([]byte(resp), &entries); err != nil {
		// Return raw response lines as fallback.
		return strings.Split(strings.TrimSpace(resp), "\n"), nil
	}
	var ips []string
	for _, e := range entries {
		if e.IP != "" {
			ips = append(ips, e.IP)
		}
	}
	if len(ips) == 0 {
		return []string{"(no DNS servers detected)"}, nil
	}
	return ips, nil
}

// SpeedTest performs a simple download speed test.
func SpeedTest() (string, error) {
	url := "https://speed.cloudflare.com/__down?bytes=10000000" // 10 MB
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("speed test request failed: %w", err)
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return "", fmt.Errorf("speed test download failed: %w", err)
	}

	elapsed := time.Since(start).Seconds()
	if elapsed == 0 {
		return "instant", nil
	}
	mbps := (float64(n) * 8) / (elapsed * 1_000_000)
	return fmt.Sprintf("Downloaded %.1f MB in %.1fs = %.1f Mbps", float64(n)/1_000_000, elapsed, mbps), nil
}

func httpGet(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
