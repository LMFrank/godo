## godo
基于golang cobra包的命令行运维工具

### 安装
下载源码，使用`make build`命令，编译对应机器系统的可执行文件

### 命令介绍：
```shell
Usage:
  godo [command]

Available Commands:
  help        Help about any command
  ping        Ping a specified IP or multiple IPs from a YAML file.
  set         Set YUM source for CentOS systems.

Flags:
  -h, --help     help for godo
```

### ping
ping命令用于测试IP是否可达，支持单IP和YAML文件批量测试
示例：
1. 测试单个IP：
`godo ping 8.8.8.8`
2. 测试多个IP，将需要测试的IP写入YAML文件，如host.yaml：
host.yaml:
```yaml
hosts:
  - 8.8.8.8
  - 114.114.114.114
  - 223.5.5.5
  - 180.76.76.76
```

### set
set命令用于设置YUM源，支持CentOS系统，默认为阿里云源
示例：
会根据系统自动选择对应的YUM源，支持CentOS6/7/8
`godo set yum`
