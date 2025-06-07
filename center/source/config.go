package source

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ListenConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// CenterConfig 定义 center 的配置结构体
// 支持多个 node
type CenterConfig struct {
	Listen ListenConfig `yaml:"listen"`
	SMTP   SMTPConfig   `yaml:"smtp"`
	Nodes  []NodeConfig `yaml:"nodes"`
}

type NodeConfig struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

// AppConfig 保存全局配置
var AppConfig *CenterConfig

// LoadConfig 解析 YAML 配置文件
func LoadConfig(path string) (*CenterConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg CenterConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	AppConfig = &cfg
	return &cfg, nil
}

// GenerateDefaultConfig 自动生成一个默认配置文件
func GenerateDefaultConfig(path string) error {
	defaultCfg := CenterConfig{
		Listen: ListenConfig{
			IP:   "0.0.0.0",
			Port: 8081,
		},
		SMTP: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     465,
			Username: "user@example.com",
			Password: "your_password",
		},
		Nodes: []NodeConfig{
			{
				Name:    "node1",
				Address: "http://192.168.1.10:8080",
				Token:   "node1_token",
			},
			{
				Name:    "node2",
				Address: "http://192.168.1.11:8080",
				Token:   "node2_token",
			},
		},
	}
	data, err := yaml.Marshal(&defaultCfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
