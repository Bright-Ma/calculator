package calculator

import (
	"encoding/xml"
	"log"
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
	log.Printf("开始加载配置文件: %s", filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("错误: 读取配置文件失败: %v", err)
		return nil, err
	}

	log.Printf("成功读取配置文件, 大小: %d 字节", len(data))

	var config Config
	err = xml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("错误: 解析XML配置失败: %v", err)
		return nil, err
	}

	log.Printf("成功解析XML配置, 包含 %d 个操作定义", len(config.Operations))

	// 记录已配置的操作
	for i, op := range config.Operations {
		status := "禁用"
		if op.Enabled {
			status = "启用"
		}
		log.Printf("操作 #%d: 类型=%s, 名称=%s, 状态=%s", i+1, op.Type, op.Name, status)
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
