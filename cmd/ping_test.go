package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockExecutor is a mock implementation of PingExecutor
type MockExecutor struct {
	output string
	err    error
}

func (m *MockExecutor) Command(name string, arg ...string) *exec.Cmd {
	cmd := &exec.Cmd{
		Path:   name,
		Args:   append([]string{name}, arg...),
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}
	cmd.Stdout.(*bytes.Buffer).WriteString(m.output)
	cmd.Stderr.(*bytes.Buffer).WriteString(m.err.Error())
	return cmd
}

// TestPingSingleIP tests the ping command for a single IP address
func TestPingSingleIP(t *testing.T) {
	// Mock the ping command output
	_ = &MockExecutor{
		output: "64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms\n",
		err:    nil,
	}

	// Run the ping command
	buf := new(bytes.Buffer)
	cmd := pingCmd
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"8.8.8.8"})
	err := cmd.Execute()

	assert.NoError(t, err, "Expected no error running the ping command")

	// Check the output format
	output := buf.String()
	assert.Contains(t, output, "64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms", "Expected specific output format")
}

// TestPingMultipleIPs tests the ping command for multiple IP addresses from a YAML file
func TestPingMultipleIPs(t *testing.T) {
	// Create a temporary YAML file
	ymlContent := `hosts:
  - 8.8.8.8
  - 1.1.1.1
`
	tempFile, err := ioutil.TempFile("", "hosts.yml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(ymlContent)
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Mock the ping command output
	_ = &MockExecutor{
		output: "64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms\n",
		err:    nil,
	}

	// Run the ping command with the config file
	cmd := pingCmd
	cmd.SetArgs([]string{"-c", tempFile.Name()})
	err = cmd.Execute()

	assert.NoError(t, err, "Expected no error running the ping command")

	// Check the response.txt file
	responseFile := "response.txt"
	data, err := ioutil.ReadFile(responseFile)
	if err != nil {
		t.Fatalf("Failed to read response file: %v", err)
	}

	// Check the output format
	output := string(data)
	assert.Contains(t, output, "8.8.8.8,64 bytes from 8.8.8.8: icmp_seq=1 ttl=56 time=12.3 ms", "Expected specific output format for 8.8.8.8")
	assert.Contains(t, output, "1.1.1.1,64 bytes from 1.1.1.1: icmp_seq=1 ttl=56 time=12.3 ms", "Expected specific output format for 1.1.1.1")
}
