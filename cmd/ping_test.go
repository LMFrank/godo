package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

// MockExecCommand is a mock function for exec.Command
var MockExecCommand = func(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// TestPingIP tests the pingIP function
func TestPingIP(t *testing.T) {
	// Mock the exec.Command function to return a known output
	mockOutput := "64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms\n64 bytes from 8.8.8.8: icmp_seq=2 ttl=56 time=11.9 ms\n"
	MockExecCommand = func(name string, arg ...string) *exec.Cmd {
		cmd := exec.Command(name, arg...)
		cmd.Stdout = bytes.NewBufferString(mockOutput)
		cmd.Stderr = bytes.NewBufferString("")
		return cmd
	}

	// Test a valid IP
	ip := "8.8.8.8"
	pingIP(ip)

	// Test an invalid IP
	ip = "invalid_ip"
	pingIP(ip)
}

// TestPingMultipleIPs tests the pingMultipleIPs function
func TestPingMultipleIPs(t *testing.T) {
	// Create a temporary YAML file with multiple IPs
	tempFile, err := os.CreateTemp("", "hosts.yml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	hostsContent := `hosts:
  - 8.8.8.8
  - 223.5.5.5
  - 180.76.76.76
`
	if _, err := tempFile.WriteString(hostsContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Mock the exec.Command function to return known outputs
	mockOutputs := map[string]string{
		"8.8.8.8":      "64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms\n64 bytes from 8.8.8.8: icmp_seq=2 ttl=56 time=11.9 ms\n",
		"223.5.5.5":    "64 bytes from 223.5.5.5: icmp_seq=1 ttl=56 time=11.3 ms\n64 bytes from 223.5.5.5: icmp_seq=2 ttl=56 time=11.9 ms\n",
		"180.76.76.76": "64 bytes from 180.76.76.76: icmp_seq=1 ttl=64 time=1.3 ms\n64 bytes from 180.76.76.76: icmp_seq=2 ttl=64 time=1.9 ms\n",
	}
	MockExecCommand = func(name string, arg ...string) *exec.Cmd {
		ip := arg[len(arg)-1]
		output := mockOutputs[ip]
		cmd := exec.Command(name, arg...)
		cmd.Stdout = bytes.NewBufferString(output)
		cmd.Stderr = bytes.NewBufferString("")
		return cmd
	}

	// Call the function
	pingMultipleIPs(tempFile.Name())

	// Read the response.csv file and verify the content
	expectedContent := `IP,Result
8.8.8.8,64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms
64 bytes from 8.8.8.8: icmp_seq=2 ttl=56 time=11.9 ms
223.5.5.5,64 bytes from 223.5.5.5: icmp_seq=1 ttl=56 time=11.3 ms
64 bytes from 223.5.5.5: icmp_seq=2 ttl=56 time=11.9 ms
180.76.76.76,64 bytes from 180.76.76.76: icmp_seq=1 ttl=64 time=1.3 ms
64 bytes from 180.76.76.76: icmp_seq=2 ttl=64 time=1.9 ms
`
	actualContent, err := os.ReadFile("response.csv")
	if err != nil {
		t.Fatalf("Failed to read response.csv: %v", err)
	}

	assert.Equal(t, expectedContent, string(actualContent))
}

// TestPingCmd tests the pingCmd function
func TestPingCmd(t *testing.T) {
	// Create a new root command
	rootCmd := &cobra.Command{}
	pingCmd := NewPingCmd()
	rootCmd.AddCommand(pingCmd)

	// Test with a single IP
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"ping", "8.8.8.8"})
	MockExecCommand = func(name string, arg ...string) *exec.Cmd {
		cmd := exec.Command(name, arg...)
		cmd.Stdout = bytes.NewBufferString("64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms\n")
		cmd.Stderr = bytes.NewBufferString("")
		return cmd
	}
	err := rootCmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms")

	// Test with a config file
	tempFile, err := os.CreateTemp("", "hosts.yml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	hostsContent := `hosts:
  - 8.8.8.8
  - 223.5.5.5
  - 180.76.76.76
`
	if _, err := tempFile.WriteString(hostsContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	rootCmd.SetArgs([]string{"ping", "-c", tempFile.Name()})
	MockExecCommand = func(name string, arg ...string) *exec.Cmd {
		ip := arg[len(arg)-1]
		output := "64 bytes from " + ip + ": icmp_seq=1 ttl=56 time=12.3 ms\n"
		cmd := exec.Command(name, arg...)
		cmd.Stdout = bytes.NewBufferString(output)
		cmd.Stderr = bytes.NewBufferString("")
		return cmd
	}
	err = rootCmd.Execute()
	assert.NoError(t, err)

	// Read the response.csv file and verify the content
	expectedContent := `IP,Result
8.8.8.8,64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms
223.5.5.5,64 bytes from 223.5.5.5: icmp_seq=1 ttl=56 time=12.3 ms
180.76.76.76,64 bytes from 180.76.76.76: icmp_seq=1 ttl=56 time=12.3 ms
`
	actualContent, err := os.ReadFile("response.csv")
	if err != nil {
		t.Fatalf("Failed to read response.csv: %v", err)
	}

	assert.Equal(t, expectedContent, string(actualContent))
}

// NewPingCmd creates a new pingCmd
func NewPingCmd() *cobra.Command {
	var configFile string

	pingCmd := &cobra.Command{
		Use:   "ping [ip | -c hosts.yml]",
		Short: "Ping a specified IP address or multiple IPs from a YAML file",
		Long: `Ping a specified IP address or multiple IPs from a YAML file and print the results to the console or a file.
For example:
godo ping 8.8.8.8
godo ping -c hosts.yml`,
		Args: cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			if configFile != "" {
				pingMultipleIPs(configFile)
			} else if len(args) > 0 {
				ip := args[0]
				pingIP(ip)
			} else {
				fmt.Println("Please provide an IP address or a config file.")
			}
		},
	}

	pingCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the YAML config file containing multiple IPs")

	return pingCmd
}
