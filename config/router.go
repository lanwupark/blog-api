package config

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	r      = NewRouter()
	routes = []*Route{} // 初始化子路由
	server *http.Server
)

// Router 小小路由封装
type Router struct {
	handlers []HTTPRequestHandler
	router   *mux.Router
}

// Route 路由 wrapper
type Route struct {
	Method          string               //方法类型
	Path            string               //路由路径
	Handler         http.HandlerFunc     //处理器
	MiddlewareFuncs []mux.MiddlewareFunc //中间件
}

// HTTPRequestHandler http请求的handler
type HTTPRequestHandler interface {
	GetRoutes() []*Route
}

func init() {
	// 注册该服务
	c.RegisterService(r)
}

// NewRouter 新建router
func NewRouter() *Router {
	return &Router{
		// 使用gorilla的mux HTTP多路复用器 它实现了http.Hanlder接口所以和http.ServeMux兼容
		router:   mux.NewRouter(),
		handlers: []HTTPRequestHandler{},
	}
}

// Config 路由配置
func (r *Router) Config(configs *Configs) {
	router := r.router
	// 配置路由
	r.configAllRoute()

	// CORS 跨域资源访问
	corsHanlder := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	server = &http.Server{
		Addr:         configs.BindAdreess,
		Handler:      corsHanlder(httplog.New(router)), //跨域访问+ httplog中间件
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	r.router = router

	log.Infof("starting server on: %s", configs.BindAdreess)
	// block
	err := server.ListenAndServe()
	if err != nil {
		log.WithError(err).Warn("shutdown")
		os.Exit(1)
	}
}

// Shutdown 实现Service接口
func (r *Router) Shutdown() {
	// 截止时间是现在+设置的绝对时间
	tx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// 如果没有处理程序 则正常关闭 如果30s后任然有请求发生 则强制关闭
	server.Shutdown(tx)
}

// 配置所有路由
func (r *Router) configAllRoute() {
	router := r.router
	for _, handler := range r.handlers {
		for _, route := range handler.GetRoutes() {
			r := router.Methods(route.Method).Subrouter()
			// 使用中间件
			r.Use(route.MiddlewareFuncs...)
			// 映射处理函数
			r.HandleFunc(route.Path, route.Handler)
		}
	}
}

// AddHTTPRequestHanlder 添加
func (r *Router) AddHTTPRequestHanlder(hanlder HTTPRequestHandler) {
	r.handlers = append(r.handlers, hanlder)
}

// GetDefaultRouter 获取路由
func GetDefaultRouter() *Router {
	return r
}
