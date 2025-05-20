package calculator

import (
	"encoding/xml"
	"os"
)

// Config 计算器配置结构
type Config struct {
	XMLName    xml.Name          `xml:"calculator"`
	Operations []OperationConfig `xml:"operations>operation"`
}

// OperationConfig 操作配置结构
type OperationConfig struct {
	Type        string `xml:"type,attr"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Enabled     bool   `xml:"enabled"`
}

// LoadConfig 从XML文件加载配置
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = xml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// GetEnabledOperations 获取所有启用的操作配置
func (c *Config) GetEnabledOperations() []OperationConfig {
	enabled := make([]OperationConfig, 0)
	for _, op := range c.Operations {
		if op.Enabled {
			enabled = append(enabled, op)
		}
	}
	return enabled
}
