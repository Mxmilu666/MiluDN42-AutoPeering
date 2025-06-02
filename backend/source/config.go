package source

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// 服务器监听配置
	Server struct {
		IP   string `yaml:"ip"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	// 认证配置
	Token string `yaml:"token"`

	// Bird配置
	Bird struct {
		PeerConfPath string `yaml:"peerconf_path"`
	} `yaml:"bird"`

	// Wireguard配置
	Wireguard struct {
		ConfigPath string `yaml:"config_path"`
	} `yaml:"wireguard"`

	// DN42网络配置
	DN42 struct {
		IPv4 string `yaml:"ipv4"`
		IPv6 string `yaml:"ipv6"`
		ASN  string `yaml:"asn"`
	} `yaml:"dn42"`
}

// 全局变量保存配置
var AppConfig *Config

// DefaultConfig 返回一个带有默认值的配置实例
func DefaultConfig() Config {
	config := Config{}
	// 初始化服务器配置
	config.Server.IP = "0.0.0.0"
	config.Server.Port = 8080

	// 认证配置
	config.Token = "change_me"

	// Bird配置
	config.Bird.PeerConfPath = "/etc/bird/peers/"

	// Wireguard配置
	config.Wireguard.ConfigPath = "/etc/wireguard/"

	// DN42网络配置
	config.DN42.IPv4 = ""
	config.DN42.IPv6 = ""
	config.DN42.ASN = ""

	return config
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// 如果文件不存在，创建默认配置文件
		if os.IsNotExist(err) {
			config := DefaultConfig()
			return createDefaultConfig(configPath, config)
		}
		return nil, fmt.Errorf("failed to read config file: %v", err)
	} // 解析配置文件
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// 配置文件验证
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	// 设置全局配置变量
	AppConfig = &config

	return AppConfig, nil
}

// createDefaultConfig 创建一个包含默认配置的新配置文件
func createDefaultConfig(configPath string, config Config) (*Config, error) {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize default config: %v", err)
	}

	// 直接使用配置数据
	finalContent := data

	if err := os.WriteFile(configPath, finalContent, 0644); err != nil {
		return nil, fmt.Errorf("failed to write default config file: %v", err)
	}

	fmt.Printf("Default config file created: %s\n", configPath)
	fmt.Println("Exiting after generating initial configuration file")
	os.Exit(0)
	return nil, nil
}

// validateConfig 验证配置是否有效
func validateConfig(config *Config) error {
	// 确保必要的字段已填写
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Token == "" || config.Token == "change_me" {
		return fmt.Errorf("please set a valid token")
	}
	if config.Bird.PeerConfPath == "" {
		return fmt.Errorf("bird peer path cannot be empty")
	}

	if config.Wireguard.ConfigPath == "" {
		return fmt.Errorf("wireguard config path cannot be empty")
	}

	// DN42网络配置验证
	if config.DN42.IPv4 == "" && config.DN42.IPv6 == "" {
		return fmt.Errorf("at least one DN42 IP address (IPv4 or IPv6) is required")
	}

	if config.DN42.ASN == "" {
		return fmt.Errorf("DN42 ASN cannot be empty")
	}

	return nil
}

// SaveConfig 将配置保存到文件
func SaveConfig(configPath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %v", err)
	}

	// 直接使用配置数据
	finalContent := data

	if err := os.WriteFile(configPath, finalContent, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
