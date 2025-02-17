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

type PingConfig struct {
	Hosts []string `yaml:"hosts"`
}

func PingIP(ip string) {
	if net.ParseIP(ip) == nil {
		fmt.Printf("Invalid IP address: %s\n", ip)
		return
	}

	command := fmt.Sprintf("ping -c 4 %s", ip)
	executor := util.DefaultCommandExecutor{}
	err, result := executor.Execute(command)
	if err != nil {
		fmt.Printf("Ping command failed with error: %s: %v\n", ip, err)
		return
	}

	fmt.Printf("Ping results for %s:\n%s", ip, result)
}

func PingMultipleIPs(configFile string) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	var config PingConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return
	}

	timestamp := time.Now().Format("20060102_15-04-05")
	dir := "./response"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	filename := fmt.Sprintf("%s/response_%s.csv", dir, timestamp)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating response file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"IP", "Result"})

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

func PingSingleIP(ip string) string {
	if net.ParseIP(ip) == nil {
		return fmt.Sprintf("Invalid IP address: %s\n", ip)
	}

	command := fmt.Sprintf("ping -c 4 %s", ip)
	executor := util.DefaultCommandExecutor{}
	err, result := executor.Execute(command)
	if err != nil {
		return fmt.Sprintf("Error pinging %s: %v\n", ip, err)
	}

	return result
}
