package config

import "github.com/apex/log"

var banner string = `                                  

██╗      █████╗ ███╗   ██╗██╗    ██╗██╗   ██╗██████╗  █████╗ ██████╗ ██╗  ██╗    ██████╗ ██╗      ██████╗  ██████╗ 
██║     ██╔══██╗████╗  ██║██║    ██║██║   ██║██╔══██╗██╔══██╗██╔══██╗██║ ██╔╝    ██╔══██╗██║     ██╔═══██╗██╔════╝ 
██║     ███████║██╔██╗ ██║██║ █╗ ██║██║   ██║██████╔╝███████║██████╔╝█████╔╝     ██████╔╝██║     ██║   ██║██║  ███╗
██║     ██╔══██║██║╚██╗██║██║███╗██║██║   ██║██╔═══╝ ██╔══██║██╔══██╗██╔═██╗     ██╔══██╗██║     ██║   ██║██║   ██║
███████╗██║  ██║██║ ╚████║╚███╔███╔╝╚██████╔╝██║     ██║  ██║██║  ██║██║  ██╗    ██████╔╝███████╗╚██████╔╝╚██████╔╝
╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚══╝╚══╝  ╚═════╝ ╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝ 
                                                                                                                   
`
var (
	c          Configs //全局配置结构体
	configured bool    //是否已配置
	finished   bool
)

// ApplicationProperties 全局配置字符串
type ApplicationProperties struct {
	DSN               string //数据库连接串
	MongoURL          string //mongodb 连接串
	MongoDatabase     string //mongodb 数据库
	RedisURL          string //redis 连接串
	BindAdreess       string //restful url
	FileBaseDir       string //文件存储基础位置
	FileMaxSize       int    //文件最大大小
	OAuthClientID     string //github OAuth 的client ID
	OAuthClientSecret string //github OAuth 的client secret
}

// Configs 所有配置
type Configs struct {
	ApplicationProperties
	services []Service
}

// Service 配置接口 实现该接口 在程序启动后会调用Service()方法
type Service interface {
	Config(c *Configs)
	Shutdown()
}

// RegisterServices 注册所有服务 扩展的服务需要写在这
func (c *Configs) RegisterServices() {
	c.RegisterService(GetDefaultRouter())
	c.RegisterService(GetConnection())
}

// LoadConfigs 批量配置所有
func (c *Configs) LoadConfigs() {
	if !configured {
		// 创建n个缓冲通道 来实现waitGroup
		done := make(chan struct{}, len(c.services))
		for _, v := range c.services {
			go func(config Service) {
				log.Infof("configuring:\t%T", config)
				config.Config(c)
				done <- struct{}{}
			}(v)
		}
		<-done
	}
	configured = true
}

// ShutdownAll 结束所有服务
func (c *Configs) ShutdownAll() {
	if !finished {
		for _, v := range c.services {
			v.Shutdown()
			log.Infof("%T\t over", v)
		}
	}
	finished = true
}

// RegisterService 注册配置
func (c *Configs) RegisterService(service Service) {
	c.services = append(c.services, service)
}

// LogBanner 打印banner
func (c *Configs) LogBanner() {
	log.Info(banner)
}

// GetConfigs 注册所有配置
func GetConfigs() *Configs {
	return &c
}
