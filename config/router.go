package config

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	r      Router
	routes []*Route
)

// Router 小小路由封装
type Router struct {
	router mux.Router
}

// Route 路由 wrapper
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func init() {
	// 初始化子路由
	routes = []*Route{}

	// 添加子路由测试
	route := &Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte("Hello Wolrd"))
		},
	}
	AddSubRoute(route)

	c.RegisterConfig(r)
}

// Config 路由配置
func (r Router) Config(configs *Configs) {
	// 使用gorilla的mux HTTP多路复用器 它实现了http.Hanlder接口所以和http.ServeMux兼容
	router := mux.NewRouter()
	routeConfig(router)

	// CORS 跨域资源访问
	corsHanlder := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	server := &http.Server{
		Addr:         configs.BindAdreess,
		Handler:      corsHanlder(httplog.New(router)), //跨域访问+ http log中间件
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		log.Infof("starting server on: %s", configs.BindAdreess)
		// block
		err := server.ListenAndServe()
		if err != nil {
			log.WithError(err).Info("shutdown")
			os.Exit(1)
		}
	}()

	r.router = *router
	// 添加钩子函数
	hookFunc(server)

}

func hookFunc(server *http.Server) {
	// 使用os/signal包里 通知某种信号来告知程序关闭服务器
	sigChannel := make(chan os.Signal)
	// 当收到终止或kill命令时，会向sigChan发送
	signal.Notify(sigChannel, os.Interrupt, os.Kill)
	// 在未收到信号前 这里是阻塞的
	sig := <-sigChannel
	log.Warnf("received terminated signal , graceful shutdown:%v", sig)
	// 截止时间是现在+设置的绝对时间
	tx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// 如果没有处理程序 则正常关闭 如果30s后任然有请求发生 则强制关闭
	server.Shutdown(tx)
}

// 配置所有路由
func routeConfig(router *mux.Router) {
	for _, route := range routes {
		r := router.Methods(route.Method).Subrouter()
		r.HandleFunc(route.Path, route.Handler)
	}
}

// AddSubRoute 添加子路由
func AddSubRoute(route *Route) *Route {
	routes = append(routes, route)
	return route
}

// AddSubRoute 添加子路由
func (Route) AddSubRoute(route *Route) *Route {
	routes = append(routes, route)
	return route
}

// GetRouter 获取路由
func GetRouter() *mux.Router {
	return &r.router
}
