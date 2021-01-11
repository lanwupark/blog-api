package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/handler"
	"github.com/lanwupark/blog-api/service"
)

var (
	c = config.GetConfigs()
)

func init() {
	flag.StringVar(&c.DSN, "dsn", "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local", "dabatase connection string")
	flag.StringVar(&c.BindAdreess, "address", ":8080", "server address")
	flag.StringVar(&c.MongoURL, "mongo-url", "mongodb://localhost:27017", "the mongo db connection string")
	flag.StringVar(&c.RedisURL, "redis-url", "redis://localhost:6379/0", "the redis db connection string")
	flag.StringVar(&c.MongoDatabase, "mongo-db", "blog", "the default mongodb database")
	flag.StringVar(&c.OAuthClientID, "oauth-id", "0fbcd40d4a1e12920596", "github oauth client id")
	flag.StringVar(&c.OAuthClientSecret, "oauth-secret", "e449dfcd57e98903e50a4ee5208475ee0009c39a", "github oauth client secret")
	flag.StringVar(&c.FileBaseDir, "file-base-dir", "files", "file base dir")
	flag.IntVar(&c.FileMaxSize, "file-max-size", 1024*1024*4, "file max size")
}

func main() {
	flag.Parse()                  // 解析参数
	log.SetLevel(log.DebugLevel)  // 设置日志级别
	registerHTTPRequestHanlders() // 先向路由服务注册路由
	c.RegisterServices()          // 注册所有服务配置
	c.LoadConfigs()               // 加载所有服务配置
	c.LogBanner()                 // 打印 banner
	tickerFunc()                  // 定时函数
	hookFunc()                    // 钩子函数
	select {}                     // 让main函数阻塞 防止程序退出
}

func registerHTTPRequestHanlders() {
	// 获取封装得默认路由 默认路由是mux的封装 它会自动加到配置里去
	router := config.GetDefaultRouter()
	// userHandler
	router.AddHTTPRequestHanlder(handler.NewUserHandler())
	// oauthHandler
	router.AddHTTPRequestHanlder(handler.NewOAuthHandler())
	router.AddHTTPRequestHanlder(handler.NewArticleHandler())
	router.AddHTTPRequestHanlder(handler.NewCommonHandler())
	// 相册 handler
	router.AddHTTPRequestHanlder(handler.NewAlbumHander())
	// 管理员
	router.AddHTTPRequestHanlder(handler.NewAdminHandler())
}

// tickerFunc 定时函数 每小时重新设置一次排行
func tickerFunc() {
	service.CalculateSort()
	ticker := time.NewTicker(1 * time.Hour)
	go func(ticker *time.Ticker) {
		for {
			<-ticker.C
			log.Info("recalculate sort")
			service.CalculateSort()
		}
	}(ticker)
}

// hookFunc 用于平滑退出程序
func hookFunc() {
	go func() {
		// 使用os/signal包里 通知某种信号来告知程序关闭服务器
		sigChannel := make(chan os.Signal)
		// 当收到终止或kill命令时，会向sigChan发送
		signal.Notify(sigChannel, os.Interrupt, os.Kill)
		// 在未收到信号前 这里是阻塞的
		sig := <-sigChannel
		log.Warnf("received terminated signal , graceful shutdown:%v", sig)
		// 结束所有服务
		c.ShutdownAll()
	}()
}
