## godo
基于golang cobra包的命令行运维工具

### 安装
```bash
git clone https://github.com/yourusername/godo.git
cd godo
make build
```

### 命令介绍：
```shell
Usage:
  godo [command]

Available Commands:
  help        Help about any command.
  net         Network-related commands.
  set         Set YUM source for CentOS systems.

Flags:
  -h, --help     help for godo
```

### net
#### ping
ping命令用于测试IP是否可达，支持单IP和YAML文件批量测试
示例：
1. 测试单个IP：
`godo net ping 8.8.8.8`
2. 测试多个IP，将需要测试的IP写入YAML文件，如host.yaml：
`godo net ping -f hosts.yaml`
host.yaml:
```yaml
hosts:
  - 8.8.8.8
  - 114.114.114.114
  - 223.5.5.5
  - 180.76.76.76
```

#### dns
聚合多个公共 DNS 服务器（如 8.8.8.8、114.114.114.114）解析同一域名，对比响应时间及结果，检测 DNS 劫持或污染。
示例：
`godo net dns -f hosts.yaml`
host.yaml:
```yaml
hosts:
  - 8.8.8.8
  - 114.114.114.114
  - 223.5.5.5
  - 180.76.76.76
```

### set
#### yum
yum命令用于设置YUM源，支持CentOS系统，默认为阿里云源
示例：
会根据系统自动选择对应的YUM源，支持CentOS6/7/8
`godo set yum`

### 项目结构
```
godo/
├── cmd/            # 命令行实现
│   ├── net/        # 网络相关的命令实现
│   │   ├── net.go  # net命令逻辑
│   ├── set/        # 设置相关的命令实现
│   │   ├── set.go  # set命令逻辑
├── pkg/            # 包含可重用的功能模块
│   ├── net/        # 网络相关的功能模块
│   └── set/        # 设置相关的功能模块
├── response/       # 可能包含响应处理相关的代码（当前为空）
├── util/           # 工具函数
│   └── command_executor.go # 命令执行器
├── go.mod          # 依赖管理
├── go.sum          # 依赖校验和
├── hosts.yaml      # 主机配置文件
├── LICENSE         # 许可证文件
├── main.go         # 程序入口
└── README.md       # 项目文档
```
