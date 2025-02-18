package util

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 定义日志配置项
type LogConfig struct {
	LogDir          string
	DirMode         os.FileMode
	MaxSize         int  // 每个日志文件的最大大小（MB）
	MaxBackups      int  // 保留的旧日志文件最大数量
	MaxAge          int  // 保留的旧日志文件最大天数
	Compress        bool // 是否压缩旧日志文件
	TimestampFormat string
}

// Logger 定义了日志记录器的接口
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
}

// DefaultLogger 实现了 Logger 接口
type DefaultLogger struct {
	logger *logrus.Logger
	cmd    string
}

// NewLogger 创建一个新的日志记录器
func NewLogger(cmd string) (*DefaultLogger, error) {
	return NewLoggerWithConfig(cmd, &LogConfig{
		LogDir:          "logs",
		DirMode:         0755,
		MaxSize:         5,
		MaxBackups:      7,
		MaxAge:          0, // 0 表示不删除旧日志
		Compress:        true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// NewLoggerWithConfig 使用自定义配置创建日志记录器
func NewLoggerWithConfig(cmd string, config *LogConfig) (*DefaultLogger, error) {
	// 创建日志目录
	if err := os.MkdirAll(config.LogDir, config.DirMode); err != nil {
		return nil, fmt.Errorf("创建日志目录失败 [cmd=%s, time=%s]: %v",
			cmd, time.Now().Format(time.RFC3339), err)
	}

	// 配置日志轮转
	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(config.LogDir, fmt.Sprintf("%s.log", timestamp))
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// 创建 logrus 日志记录器
	logger := logrus.New()
	logger.SetOutput(lumberJackLogger)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: config.TimestampFormat,
		FullTimestamp:   true,
	})

	// 记录初始化日志
	logger.WithField("command", cmd).Info("日志记录器初始化成功")

	return &DefaultLogger{
		logger: logger,
		cmd:    cmd,
	}, nil
}

// Close 关闭日志文件（由于使用了 lumberjack，不再需要手动关闭文件）
func (l *DefaultLogger) Close() error {
	return nil
}

// Info 记录信息级别的日志
func (l *DefaultLogger) Info(format string, v ...interface{}) {
	l.logger.WithField("command", l.cmd).Infof(format, v...)
}

// Error 记录错误级别的日志
func (l *DefaultLogger) Error(format string, v ...interface{}) {
	l.logger.WithField("command", l.cmd).Errorf(format, v...)
}

// Debug 记录调试级别的日志
func (l *DefaultLogger) Debug(format string, v ...interface{}) {
	l.logger.WithField("command", l.cmd).Debugf(format, v...)
}

// Warn 记录警告级别的日志
func (l *DefaultLogger) Warn(format string, v ...interface{}) {
	l.logger.WithField("command", l.cmd).Warnf(format, v...)
}
