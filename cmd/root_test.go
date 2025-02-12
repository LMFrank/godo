package cmd

import (
	"bytes"
	"testing"

	_ "github.com/spf13/cobra"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "invalid command",
			args:    []string{"invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)
			if err := rootCmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRootCommandFlags(t *testing.T) {
	cmd := rootCmd

	if cmd.Flags().Lookup("toggle") == nil {
		t.Error("toggle flag not found")
	}
}

func TestRootCommandHelpOutput(t *testing.T) {
	cmd := rootCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := b.String()
	if len(output) == 0 {
		t.Error("Expected non-empty help output")
	}

	if !contains(output, "Usage:") {
		t.Error("Expected usage section in help output")
	}
	if !contains(output, "Available Commands:") {
		t.Error("Expected commands section in help output")
	}
	if !contains(output, "Flags:") {
		t.Error("Expected flags section in help output")
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
