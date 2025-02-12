package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ExecuteCmd(command string) (err error, result string) {
	cmd := exec.Command("/bin/bash", "-c", command)
	fmt.Printf("[shell] Executing command: %s\n", command)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err, ""
	}

	return nil, out.String()
}
