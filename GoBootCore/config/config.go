package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config 定义了应用程序的配置结构，包含服务器、数据库、JWT 和 RabbitMQ 的配置信息。
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Type     string `yaml:"type"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	JWT struct {
		Secret                 string `yaml:"secret"`
		ExpirationMilliseconds int    `yaml:"expiration_milliseconds"`
	} `yaml:"jwt"`
	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"rabbitmq"`
	Sqlite struct {
		TileDataPath string `yaml:"tile_data_path"`
	} `yaml:"sqlite"`
}

// LoadConfig 从配置文件中加载配置信息。
// 该函数会优先尝试从可执行文件所在目录加载 config.yaml，
// 如果未找到，则尝试从当前工作目录加载。
//
// 返回值：
//   - *Config: 成功时返回解析后的配置对象
//   - error: 失败时返回错误信息
func LoadConfig() (*Config, error) {
	config := &Config{}

	// 获取可执行文件所在的目录
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// 定义可能的配置文件路径
	exeConfigPath := filepath.Join(exeDir, "config.yaml")
	rootConfigPath := "config.yaml" // 假设项目根目录为当前工作目录

	var fileContent []byte
	var configPath string

	// 优先尝试从可执行文件所在目录读取配置文件
	if _, err := os.Stat(exeConfigPath); err == nil {
		fileContent, err = os.ReadFile(exeConfigPath)
		configPath = exeConfigPath
	} else if os.IsNotExist(err) {
		// 若可执行文件目录下不存在，则尝试从项目根目录读取
		fileContent, err = os.ReadFile(rootConfigPath)
		configPath = rootConfigPath
	} else {
		return nil, fmt.Errorf("failed to check config file existence: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// 将 YAML 格式的配置内容解析到 Config 结构体中
	err = yaml.Unmarshal(fileContent, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from %s: %w", configPath, err)
	}

	return config, nil
}
