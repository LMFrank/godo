package util

// ExecuteCmd 使用 DefaultCommandExecutor 执行命令
func ExecuteCmd(command string) (err error, result string) {
	executor := &DefaultCommandExecutor{}
	return executor.Execute(command)
}
