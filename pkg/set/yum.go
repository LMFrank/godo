package set

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// SetCentOSYumSource 将 CentOS 的 YUM 源更改为阿里云源。
func SetCentOSYumSource() {
	// 检测 CentOS 版本
	version, err := DetectCentOSVersion()
	if err != nil {
		logrus.Errorf("检测 CentOS 版本时出错: %v", err)
		return
	}
	logrus.Infof("检测到 CentOS 版本: %s", version)

	// 备份原始的 CentOS-Base.repo 文件
	logrus.Info("开始备份原始 CentOS-Base.repo 文件")
	backupCmd := exec.Command("mv", "/etc/yum.repos.d/CentOS-Base.repo",
		"/etc/yum.repos.d/CentOS-Base.repo.backup")
	if err := backupCmd.Run(); err != nil {
		logrus.Errorf("备份 CentOS-Base.repo 时出错: %v", err)
		return
	}
	logrus.Info("成功备份 CentOS-Base.repo 文件")

	// 根据不同的 CentOS 版本选择对应的阿里云源 URL
	var url string
	switch version {
	case "6":
		url = "https://mirrors.aliyun.com/repo/Centos-vault-6.10.repo"
	case "7":
		url = "https://mirrors.aliyun.com/repo/Centos-7.repo"
	case "8":
		url = "https://mirrors.aliyun.com/repo/Centos-vault-8.5.2111.repo"
	default:
		logrus.Errorf("不支持的 CentOS 版本: %s", version)
		return
	}
	logrus.Infof("选择阿里云源 URL: %s", url)

	// 下载新的 CentOS-Base.repo 文件
	logrus.Info("开始下载新的 CentOS-Base.repo 文件")
	downloadCmd := exec.Command("wget", "-O", "/etc/yum.repos.d/CentOS-Base.repo", url)
	if err := downloadCmd.Run(); err != nil {
		logrus.Errorf("下载 CentOS-Base.repo 时出错: %v", err)
		return
	}
	logrus.Info("成功下载新的 CentOS-Base.repo 文件")

	// 更新 YUM 缓存
	logrus.Info("开始更新 YUM 缓存")
	makecachedCmd := exec.Command("yum", "makecache")
	if err := makecachedCmd.Run(); err != nil {
		logrus.Errorf("运行 yum makecache 时出错: %v", err)
		return
	}
	logrus.Info("成功更新 YUM 缓存")

	// 删除非阿里云 ECS 用户可能会遇到的无法解析主机问题的行
	logrus.Info("开始处理非阿里云 ECS 用户配置")
	setCmd := exec.Command("sed", "-i", "-e", "/mirrors.cloud.aliyuncs.com/d",
		"-e", "/mirrors.aliyuncs.com/d", "/etc/yum.repos.d/CentOS-Base.repo")
	if err := setCmd.Run(); err != nil {
		logrus.Errorf("从 CentOS-Base.repo 中删除特定行时出错: %v", err)
		return
	}
	logrus.Info("成功处理非阿里云 ECS 用户配置")

	// 输出成功信息
	logrus.Infof("CentOS %s YUM 源已成功更改为阿里云", version)
}

// DetectCentOSVersion 检测当前系统的 CentOS 版本。
// 返回值：
//   - string: CentOS 版本号
//   - error: 如果有错误发生，则返回错误信息
func DetectCentOSVersion() (string, error) {
	logrus.Info("开始检测 CentOS 版本")
	cmd := exec.Command("cat", "/etc/redhat-release")
	output, err := cmd.Output()
	if err != nil {
		logrus.Errorf("读取 /etc/redhat-release 文件时出错: %v", err)
		return "", err
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
		err := fmt.Errorf("不支持的 CentOS 版本: %s", string(output))
		logrus.Error(err)
		return "", err
	}

	logrus.Infof("成功检测到 CentOS 版本: %s", version)
	return version, nil
}
