package set

import (
	"fmt"
	"os/exec"
	"strings"
)

func SetCentOSYumSource() {
	version, err := DetectCentOSVersion()
	if err != nil {
		fmt.Printf("Error detecting CentOS version: %v\n", err)
		return
	}

	backupCmd := exec.Command("mv", "/etc/yum.repos.d/CentOS-Base.repo",
		"/etc/yum.repos.d/CentOS-Base.repo.backup")
	if err := backupCmd.Run(); err != nil {
		fmt.Printf("Error backing up CentOS-Base.repo: %v\n", err)
		return
	}

	var url string
	switch version {
	case "6":
		url = "https://mirrors.aliyun.com/repo/Centos-vault-6.10.repo"
	case "7":
		url = "https://mirrors.aliyun.com/repo/Centos-7.repo"
	case "8":
		url = "https://mirrors.aliyun.com/repo/Centos-vault-8.5.2111.repo"
	default:
		fmt.Printf("Unsupported CentOS version: %s\n", version)
		return
	}

	downloadCmd := exec.Command("wget", "-O", "/etc/yum.repos.d/CentOS-Base.repo", url)
	if err := downloadCmd.Run(); err != nil {
		fmt.Printf("Error downloading CentOS-Base.repo: %v\n", err)
		return
	}

	makecachedCmd := exec.Command("yum", "makecache")
	if err := makecachedCmd.Run(); err != nil {
		fmt.Printf("Error running yum makecache: %v\n", err)
		return
	}

	// 非阿里云ECS用户会出现 Couldn't resolve host 'mirrors.cloud.aliyuncs.com' 信息，不影响使用。
	setCmd := exec.Command("sed", "-i", "-e", "/mirrors.cloud.aliyuncs.com/d",
		"-e", "/mirrors.aliyuncs.com/d", "/etc/yum.repos.d/CentOS-Base.repo")
	if err := setCmd.Run(); err != nil {
		fmt.Printf("Error removing specific lines from CentOS-Base.repo: %v\n", err)
		return
	}

	fmt.Printf("CentOS %s YUM source has been successfully changed to Aliyun.\n", version)
}

func DetectCentOSVersion() (string, error) {
	cmd := exec.Command("cat", "/etc/redhat-release")
	output, err := cmd.Output()
	if err != nil {
		return "", nil
	}

	version := ""
	switch {
	case strings.Contains(string(output), "release 6"):
		version = "6"
	case strings.Contains(string(output), "release 7"):
		version = "7"
	case strings.Contains(string(output), "release 8"):
		version = "8"
	default:
		return "", fmt.Errorf("unsupported CentOS version: %s", string(output))
	}

	return version, nil
}
