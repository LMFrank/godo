## godo
基于golang cobra包的命令行运维工具

### 安装
```bash
git clone https://github.com/yourusername/godo.git
cd godo
make build
```
## 当前功能特性

- 🚀 单IP即时ping测试
- 📁 支持YAML配置文件批量测试
- 📊 CSV格式结果输出
- 🔧 Cobra框架驱动命令行交互

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

### 项目结构
```
godo/
├── cmd/            # 命令行实现
│   ├── ping.go     # ping命令逻辑
│   └── root.go     # 根命令配置
├── util/           # 工具函数
│   └── execute.go  # 命令执行器
├── go.mod          # 依赖管理
├── main.go         # 程序入口
└── README.md       # 项目文档
```