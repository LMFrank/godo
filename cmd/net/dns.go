package net

import (
	"fmt"
	"log"
	"os"

	"github.com/LMFrank/godo/pkg/net"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "DNS 解析和响应时间检测",
	Long:  `聚合多个公共 DNS 服务器（如 8.8.8.8、114.114.114.114）解析同一域名，对比响应时间及结果，检测 DNS 劫持或污染。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: godo net dns <domain> -f <filename>")
			return
		}
		domain := args[0]

		hostsFile, _ := cmd.Flags().GetString("hosts")
		if hostsFile == "" {
			fmt.Println("请提供 hosts.yaml 文件路径")
			return
		}

		yamlFile, err := os.ReadFile(hostsFile)
		if err != nil {
			log.Fatalf("读取文件失败: %v", err)
		}

		var hosts struct {
			Servers []string `yaml:"hosts"`
		}

		err = yaml.Unmarshal(yamlFile, &hosts)
		if err != nil {
			log.Fatalf("解析 YAML 失败: %v", err)
		}

		results, err := net.ResolveDNS(domain, hosts.Servers)
		if err != nil {
			log.Fatalf("DNS 解析失败: %v", err)
		}

		fmt.Printf("DNS 解析结果 (%s):\n", domain)
		for _, result := range results {
			if result.Error != nil {
				fmt.Printf("服务器: %s, 错误: %v\n", result.Server, result.Error)
			} else {
				fmt.Printf("服务器: %s, IP: %s, 响应时间: %v\n", result.Server, result.IP, result.ResponseTime)
			}
		}
	},
}

func init() {
	dnsCmd.Flags().StringP("hosts", "f", "", "hosts.yaml 文件路径")
}
