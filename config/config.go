package config

import "github.com/apex/log"

var (
	c          Configs //全局配置结构体
	configured bool    //是否已配置
)

// GlobalConfigStrings 全局配置字符串
type GlobalConfigStrings struct {
	DSN         string //数据库连接串
	BindAdreess string //url
	FileBaseDir string //文件存储基础位置
}

// Configs 所有配置
type Configs struct {
	GlobalConfigStrings
	configs []Config
}

// Config 配置接口 实现该接口 在程序启动后会调用Config()方法
type Config interface {
	Config(c *Configs)
}

// DoConfigAll 批量配置所有
func (c *Configs) DoConfigAll() {
	if !configured {
		// 创建n个缓冲通道 来实现waitGroup
		done := make(chan struct{}, len(c.configs))
		for _, v := range c.configs {
			go func(config Config) {
				log.Infof("configuring:\t%T", config)
				config.Config(c)
				done <- struct{}{}
			}(v)
		}
		<-done
	}
	configured = true
}

// RegisterConfig 注册配置
func (c *Configs) RegisterConfig(config Config) {
	c.configs = append(c.configs, config)
}

// GetConfigs 注册所有配置
func GetConfigs() *Configs {
	return &c
}
