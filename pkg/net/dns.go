package net

import (
	"context"
	"fmt"
	"github.com/LMFrank/godo/util"
	"net"
	"time"
)

var logger *util.DefaultLogger

func init() {
	var err error
	logger, err = util.NewLogger("net dns")
	if err != nil {
		fmt.Printf("初始化日志记录器失败: %v\n", err)
	}
}

// DNSResult 保存 DNS 查询的结果。
type DNSResult struct {
	Server       string        // 用于查询的 DNS 服务器地址
	IP           string        // 解析得到的 IP 地址
	ResponseTime time.Duration // 响应时间
	Error        error         // 错误信息，如果有的话
}

// ResolveDNS 向多个 DNS 服务器查询指定域名，并返回查询结果。
// 参数：
//   - domain: 要查询的域名
//   - servers: DNS 服务器列表
//
// 返回值：
//   - []DNSResult: 每个 DNS 服务器的查询结果列表
//   - error: 如果有错误发生，则返回错误信息
func ResolveDNS(domain string, servers []string) ([]DNSResult, error) {
	logger.Info("开始解析域名 %s，使用 DNS 服务器: %v", domain, servers)
	results := make([]DNSResult, len(servers))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i, server := range servers {
		go func(i int, server string) {
			start := time.Now()
			ip, err := resolveSingleDNS(ctx, domain, server)
			duration := time.Since(start)
			results[i] = DNSResult{
				Server:       server,
				IP:           ip,
				ResponseTime: duration,
				Error:        err,
			}
		}(i, server)
	}

	// 等待所有 goroutine 完成
	time.Sleep(5 * time.Second)

	return results, nil
}

// resolveSingleDNS 向单个 DNS 服务器查询指定域名，并返回解析结果。
// 参数：
//   - ctx: 上下文，用于超时控制
//   - domain: 要查询的域名
//   - server: DNS 服务器地址
//
// 返回值：
//   - string: 解析得到的 IP 地址
//   - error: 如果有错误发生，则返回错误信息
func resolveSingleDNS(ctx context.Context, domain, server string) (string, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:53", server))
	if err != nil {
		logger.Error("连接 DNS 服务器失败: %v", err)
		return "", err
	}
	defer conn.Close()

	logger.Debug("成功连接到 DNS 服务器 %s", server)

	// 构建 DNS 查询消息（简单的 A 记录查询）
	query := []byte{
		0x00, 0x01, // 事务 ID
		0x01, 0x00, // 标志：请求递归
		0x00, 0x01, // 问题数量
		0x00, 0x00, // 回答资源记录数
		0x00, 0x00, // 权威资源记录数
		0x00, 0x00, // 额外资源记录数
	}

	// 将域名转换为 DNS 标签格式
	labels := domainToLabels(domain)
	query = append(query, labels...)

	query = append(query, []byte{
		0x00,       // 标签结束符
		0x00, 0x01, // 类型：A 记录
		0x00, 0x01, // 类：IN
	}...)

	_, err = conn.Write(query)
	if err != nil {
		logger.Error("发送 DNS 查询失败: %v", err)
		return "", err
	}

	logger.Debug("已发送 DNS 查询请求到服务器 %s", server)

	buf := make([]byte, 512)
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		logger.Error("从 DNS 服务器读取响应失败: %v", err)
		return "", err
	}

	logger.Debug("从服务器 %s 接收到响应", server)

	if n < 12 {
		return "", fmt.Errorf("无效的 DNS 响应")
	}
	if buf[2]&0x80 != 0x80 {
		return "", fmt.Errorf("DNS 响应错误")
	}

	// 检查响应码
	rcode := buf[3] & 0x0F
	if rcode != 0 {
		return "", fmt.Errorf("DNS 响应码错误: %d", rcode)
	}

	// 检查回答数量
	answerCount := (buf[6] << 8) | buf[7]
	if answerCount == 0 {
		return "", fmt.Errorf("DNS 响应中没有回答")
	}

	// 跳过头部和问题部分
	pos := 12
	for pos < n {
		if buf[pos] == 0x00 {
			pos++
			break
		}
		pos++
	}
	pos += 4 // 跳过 QTYPE 和 QCLASS

	// 解析回答记录
	for i := 0; i < int(answerCount); i++ {
		if pos+10 > n {
			return "", fmt.Errorf("无效的 DNS 响应")
		}
		if buf[pos] != 0xC0 || buf[pos+1] != 0x0C {
			return "", fmt.Errorf("无效的 DNS 响应")
		}
		pos += 10 // 跳过 NAME, TYPE, CLASS, TTL, RDLENGTH

		if pos+4 > n {
			return "", fmt.Errorf("无效的 DNS 响应")
		}
		ip := fmt.Sprintf("%d.%d.%d.%d", buf[pos], buf[pos+1], buf[pos+2], buf[pos+3])
		pos += 4

		return ip, nil
	}

	return "", fmt.Errorf("DNS 响应中没有回答")
}

// domainToLabels 将域名转换为 DNS 标签格式。
// 参数：
//   - domain: 要转换的域名
//
// 返回值：
//   - []byte: 转换后的 DNS 标签格式字节序列
func domainToLabels(domain string) []byte {
	var labels []byte
	parts := []rune(domain)
	start := 0
	for i, r := range parts {
		if r == '.' {
			labels = append(labels, byte(i-start))
			labels = append(labels, []byte(string(parts[start:i]))...)
			start = i + 1
		}
	}
	labels = append(labels, byte(len(parts)-start))
	labels = append(labels, []byte(string(parts[start:]))...)
	return labels
}
