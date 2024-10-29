package main

import (
	"fmt"
	"net"
	"os/exec"
)

func main() {
	ip := "8.8.8.8"
	ping(ip)
}

func ping(ip string) {
	// 检查 IP 地址是否有效
	if net.ParseIP(ip) == nil {
		fmt.Printf("Invalid IP address: %s\n", ip)
		return
	}

	// 执行 ping 命令
	cmd := exec.Command("ping", "-c", "4", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error pinging %s: %v\n", ip, err)
		return
	}

	fmt.Printf("Ping results for %s:\n%s", ip, output)
}
