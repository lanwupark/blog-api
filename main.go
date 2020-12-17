package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/handler"
)

var (
	c = config.GetConfigs()
)

func init() {
	flag.StringVar(&c.DSN, "dsn", "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local", "dabatase connection string")
	flag.StringVar(&c.BindAdreess, "address", ":8080", "server address")
	flag.StringVar(&c.MongoURL, "mongo-url", "mongodb://localhost:27017", "the mongo db connection string")
	flag.StringVar(&c.MongoDatabase, "mongo-db", "blog", "the default mongodb database")
	flag.StringVar(&c.OAuthClientID, "oauth-id", "", "github oauth client id")
	flag.StringVar(&c.OAuthClientID, "oauth-secret", "", "github oauth client secret")
}

func main() {
	flag.Parse()                  // 解析参数
	registerHTTPRequestHanlders() //先向路由服务注册路由
	c.RegisterServices()          // 注册所有服务配置
	c.LoadConfigs()               // 加载所有服务配置
	c.LogBanner()                 // 打印 banner
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
	router.AddHTTPRequestHanlder(handler.NewCategoryHandler())
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
