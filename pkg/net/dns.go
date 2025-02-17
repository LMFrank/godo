package net

import (
	"context"
	"fmt"
	"net"
	"time"
)

// DNSResult holds the result of a DNS query.
type DNSResult struct {
	Server       string
	IP           string
	ResponseTime time.Duration
	Error        error
}

// ResolveDNS queries multiple DNS servers for a domain and returns the results.
func ResolveDNS(domain string, servers []string) ([]DNSResult, error) {
	results := make([]DNSResult, len(servers))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i, server := range servers {
		go func(i int, server string) {
			start := time.Now()
			ip, err := resolveSingleDNS(ctx, domain, server)
			duration := time.Since(start)
			results[i] = DNSResult{
				Server:       server,
				IP:           ip,
				ResponseTime: duration,
				Error:        err,
			}
		}(i, server)
	}

	// Wait for all goroutines to finish
	time.Sleep(5 * time.Second)

	return results, nil
}

func resolveSingleDNS(ctx context.Context, domain, server string) (string, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:53", server))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// DNS query message (simple A record query)
	query := []byte{
		0x00, 0x01, // Transaction ID
		0x01, 0x00, // Flags: recursion desired
		0x00, 0x01, // Questions
		0x00, 0x00, // Answer RRs
		0x00, 0x00, // Authority RRs
		0x00, 0x00, // Additional RRs
	}

	// Convert domain to DNS label format
	labels := domainToLabels(domain)
	query = append(query, labels...)

	query = append(query, []byte{
		0x00,       // End of labels
		0x00, 0x01, // Type: A
		0x00, 0x01, // Class: IN
	}...)

	_, err = conn.Write(query)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 512)
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	if n < 12 {
		return "", fmt.Errorf("invalid DNS response")
	}
	if buf[2]&0x80 != 0x80 {
		return "", fmt.Errorf("DNS response error")
	}

	// Check response code
	rcode := buf[3] & 0x0F
	if rcode != 0 {
		return "", fmt.Errorf("DNS response code error: %d", rcode)
	}

	// Check answer count
	answerCount := (buf[6] << 8) | buf[7]
	if answerCount == 0 {
		return "", fmt.Errorf("no answer in DNS response")
	}

	// Skip header and question section
	pos := 12
	for pos < n {
		if buf[pos] == 0x00 {
			pos++
			break
		}
		pos++
	}
	pos += 4 // Skip QTYPE and QCLASS

	// Parse answer records
	for i := 0; i < int(answerCount); i++ {
		if pos+10 > n {
			return "", fmt.Errorf("invalid DNS response")
		}
		if buf[pos] != 0xC0 || buf[pos+1] != 0x0C {
			return "", fmt.Errorf("invalid DNS response")
		}
		pos += 10 // Skip NAME, TYPE, CLASS, TTL, RDLENGTH

		if pos+4 > n {
			return "", fmt.Errorf("invalid DNS response")
		}
		ip := fmt.Sprintf("%d.%d.%d.%d", buf[pos], buf[pos+1], buf[pos+2], buf[pos+3])
		pos += 4

		return ip, nil
	}

	return "", fmt.Errorf("no answer in DNS response")
}

func domainToLabels(domain string) []byte {
	var labels []byte
	parts := []rune(domain)
	start := 0
	for i, r := range parts {
		if r == '.' {
			labels = append(labels, byte(i-start))
			labels = append(labels, []byte(string(parts[start:i]))...)
			start = i + 1
		}
	}
	labels = append(labels, byte(len(parts)-start))
	labels = append(labels, []byte(string(parts[start:]))...)
	return labels
}
