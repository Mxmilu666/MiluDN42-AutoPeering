package source

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Mxmilu666/MiluDN42-AutoPeering/center/source/logger"
	"gopkg.in/yaml.v3"
)

type NodeConfig struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
}

type CenterConfig struct {
	Listen struct {
		IP   string `yaml:"ip"`
		Port int    `yaml:"port"`
	} `yaml:"listen"`

	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		From     string `yaml:"from"`
	} `yaml:"smtp"`

	Nodes []NodeConfig `yaml:"nodes"`
}

// AppConfig 保存全局配置
var AppConfig *CenterConfig

// DefaultConfig 返回一个带有默认值的配置实例
func DefaultConfig() CenterConfig {
	cfg := CenterConfig{}
	cfg.Listen.IP = "0.0.0.0"
	cfg.Listen.Port = 8081
	cfg.SMTP.Host = "smtp.example.com"
	cfg.SMTP.Port = 465
	cfg.SMTP.Username = "user@example.com"
	cfg.SMTP.Password = "your_password"
	cfg.SMTP.From = "noreply@example.com"
	cfg.Nodes = []NodeConfig{
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
	}
	return cfg
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(configPath string) (*CenterConfig, error) {
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			return createDefaultConfig(configPath, cfg)
		}
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	var cfg CenterConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	AppConfig = &cfg
	return AppConfig, nil
}

// createDefaultConfig 创建一个包含默认配置的新配置文件
func createDefaultConfig(configPath string, cfg CenterConfig) (*CenterConfig, error) {
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize default config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write default config file: %v", err)
	}
	logger.Info("Exiting after generating initial configuration file")
	os.Exit(0)
	return nil, nil
}

// validateConfig 验证配置是否有效
func validateConfig(cfg *CenterConfig) error {
	if cfg.Listen.Port <= 0 || cfg.Listen.Port > 65535 {
		return fmt.Errorf("invalid listen port: %d", cfg.Listen.Port)
	}
	if cfg.SMTP.Host == "" {
		return fmt.Errorf("smtp host cannot be empty")
	}
	if cfg.SMTP.Username == "" || cfg.SMTP.Password == "" {
		return fmt.Errorf("smtp username/password cannot be empty")
	}
	if cfg.SMTP.From == "" {
		return fmt.Errorf("smtp from cannot be empty")
	}
	if len(cfg.Nodes) == 0 {
		return fmt.Errorf("at least one node is required")
	}
	for i, n := range cfg.Nodes {
		if n.Name == "" || n.Address == "" || n.Token == "" {
			return fmt.Errorf("node[%d] fields cannot be empty", i)
		}
	}
	return nil
}

// SaveConfig 将配置保存到文件
func SaveConfig(configPath string, cfg *CenterConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	return nil
}
