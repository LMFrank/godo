package util

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CommandExecutor 定义了一个执行命令的接口
type CommandExecutor interface {
	Execute(command string) (err error, result string)
}

// DefaultCommandExecutor 实现了 CommandExecutor 接口
type DefaultCommandExecutor struct{}

// Execute 执行命令并返回结果
func (d *DefaultCommandExecutor) Execute(command string) (err error, result string) {
	cmd := exec.Command("/bin/bash", "-c", command)
	fmt.Printf("[shell] Executing command: %s\n", command)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err, ""
	}

	return nil, out.String()
}
