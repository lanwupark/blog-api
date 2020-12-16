package config

import (
	"github.com/apex/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	conn = new(Connection)
	db   *sqlx.DB
)

func init() {
	c.RegisterService(conn)
}

// Connection 数据库连接
type Connection struct {
}

// Config 实现配置接口
func (c *Connection) Config(configs *Configs) {
	db = sqlx.MustConnect("mysql", configs.DSN)
	log.Info("connect database successfully")
}

// Shutdown 结束
func (c *Connection) Shutdown() {
	db.Close()
}

// GetDB 获取DB连接
func GetDB() *sqlx.DB {
	return db
}
