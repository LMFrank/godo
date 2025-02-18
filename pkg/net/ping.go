package net

import (
	"encoding/csv"
	"fmt"
	"github.com/LMFrank/godo/util"
	"gopkg.in/yaml.v2"
	"net"
	"os"
	"sync"
	"time"
)

var pingLogger *util.DefaultLogger

func init() {
	var err error
	pingLogger, err = util.NewLogger("net ping")
	if err != nil {
		fmt.Printf("初始化日志记录器失败: %v\n", err)
	}
}

// PingConfig 定义了从配置文件中读取的 ping 配置结构。
type PingConfig struct {
	Hosts []string `yaml:"hosts"` // 要 ping 的主机列表
}

// PingIP 向指定的 IP 地址发送 ping 请求，并输出结果。
// 参数：
//   - ip: 要 ping 的 IP 地址
func PingIP(ip string) {
	if net.ParseIP(ip) == nil {
		fmt.Printf("无效的 IP 地址: %s\n", ip)
		return
	}

	command := fmt.Sprintf("ping -c 4 %s", ip)
	executor := util.DefaultCommandExecutor{}
	err, result := executor.Execute(command)
	if err != nil {
		pingLogger.Error("Ping 命令执行失败，错误信息: %s: %v", ip, err)
		return
	}

	fmt.Printf("Ping 结果为 %s:\n%s", ip, result)
}

// PingMultipleIPs 从配置文件中读取多个 IP 地址并依次进行 ping 操作，将结果保存到 CSV 文件中。
// 参数：
//   - configFile: 包含要 ping 的 IP 地址列表的配置文件路径
func PingMultipleIPs(configFile string) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		pingLogger.Error("读取配置文件时出错: %v", err)
		return
	}

	var config PingConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		pingLogger.Error("解析配置文件时出错: %v", err)
		return
	}

	timestamp := time.Now().Format("20060102_15-04-05")
	dir := "./response"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		pingLogger.Error("创建目录时出错: %v", err)
		return
	}
	filename := fmt.Sprintf("%s/response_%s.csv", dir, timestamp)
	file, err := os.Create(filename)
	if err != nil {
		pingLogger.Error("创建响应文件时出错: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"IP", "结果"})

	results := make(chan []string)
	var wg sync.WaitGroup

	for _, ip := range config.Hosts {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			result := PingSingleIP(ip)
			results <- []string{ip, result}
		}(ip)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		writer.Write(result)
	}
}

// PingSingleIP 向单个 IP 地址发送 ping 请求，并返回结果字符串。
// 参数：
//   - ip: 要 ping 的 IP 地址
//
// 返回值：
//   - string: ping 操作的结果或错误信息
func PingSingleIP(ip string) string {
	if net.ParseIP(ip) == nil {
		return fmt.Sprintf("无效的 IP 地址: %s\n", ip)
	}

	command := fmt.Sprintf("ping -c 4 %s", ip)
	executor := util.DefaultCommandExecutor{}
	err, result := executor.Execute(command)
	if err != nil {
		return fmt.Sprintf("Ping %s 时出错: %v\n", ip, err)
	}

	return result
}
