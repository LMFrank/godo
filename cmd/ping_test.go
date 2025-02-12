package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"strings"
	"testing"
)

// MockCommandExecutor is a mock implementation of util.CommandExecutor
type MockCommandExecutor struct {
	mock.Mock
}

// Execute is a mock implementation of the Execute method
func (m *MockCommandExecutor) Execute(command string) (error, string) {
	args := m.Called(command)
	return args.Error(0), args.String(1)
}

func TestPingIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{"Valid IP", "8.8.8.8", "PING 8.8.8.8 (8.8.8.8): 56 data bytes\n64 bytes from 8.8.8.8: icmp_seq=1 ttl=116 time=12.345 ms\n"},
		{"Invalid IP", "invalid-ip", "Invalid IP address: invalid-ip\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockCommandExecutor)
			if tt.ip != "invalid-ip" {
				mockExecutor.On("Execute", fmt.Sprintf("ping -c 4 %s", tt.ip)).Return(nil, tt.expected)
			}

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			defer r.Close()
			defer w.Close()

			old := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = old }()

			pingIP(tt.ip, mockExecutor)

			w.Close()
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}

			if tt.ip == "invalid-ip" {
				assert.Contains(t, buf.String(), tt.expected)
			} else {
				assert.Equal(t, fmt.Sprintf("Ping results for %s:\n%s", tt.ip, tt.expected), buf.String())
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}

func TestPingMultipleIPs(t *testing.T) {
	yamlContent := []byte(`
hosts:
  - 8.8.8.8
  - 8.8.4.4
`)
	tempFile, err := os.CreateTemp("", "hosts.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(yamlContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	tests := []struct {
		name       string
		configFile string
		expected   []string
	}{
		{"Valid Config", tempFile.Name(), []string{"8.8.8.8,PING 8.8.8.8 (8.8.8.8): 56 data bytes\n64 bytes from 8.8.8.8: icmp_seq=1 ttl=116 time=12.345 ms\n", "8.8.4.4,PING 8.8.4.4 (8.8.4.4): 56 data bytes\n64 bytes from 8.8.4.4: icmp_seq=1 ttl=116 time=12.345 ms\n"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := new(MockCommandExecutor)
			for _, ip := range tt.expected {
				split := strings.Split(ip, ",")
				mockExecutor.On("Execute", fmt.Sprintf("ping -c 4 %s", split[0])).Return(nil, split[1])
			}

			// Clear the response directory before running the test
			dir := "./response"
			if err := os.RemoveAll(dir); err != nil {
				t.Fatalf("Failed to remove response directory: %v", err)
			}
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				t.Fatalf("Failed to create response directory: %v", err)
			}

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Failed to create pipe: %v", err)
			}
			defer r.Close()
			defer w.Close()

			old := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = old }()

			pingMultipleIPs(tt.configFile, mockExecutor)

			w.Close()
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			if err != nil {
				t.Fatalf("Failed to read from pipe: %v", err)
			}

			files, err := os.ReadDir(dir)
			if err != nil {
				t.Fatalf("Failed to read response directory: %v", err)
			}
			if len(files) != 1 {
				t.Fatalf("Expected 1 file in response directory, got %d", len(files))
			}

			filename := fmt.Sprintf("%s/%s", dir, files[0].Name())
			file, err := os.Open(filename)
			if err != nil {
				t.Fatalf("Failed to open response file: %v", err)
			}
			defer file.Close()

			reader := csv.NewReader(file)
			records, err := reader.ReadAll()
			if err != nil {
				t.Fatalf("Failed to read CSV file: %v", err)
			}

			// Skip the header row if it exists
			if len(records) > 0 && records[0][0] == "IP" {
				records = records[1:]
			}

			for i, record := range records {
				expected := strings.Split(tt.expected[i], ",")
				assert.Equal(t, expected[0], record[0])
				assert.Contains(t, record[1], expected[1])
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}
