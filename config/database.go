package config

import (
	"github.com/apex/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	conn Connection
)

func init() {
	c.RegisterConfig(conn)
}

// Connection 数据库连接
type Connection struct {
	DB *gorm.DB
}

// Config 实现配置接口
func (c Connection) Config(configs *Configs) {
	var err error
	conn.DB, err = gorm.Open(mysql.Open(configs.DSN), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	log.Info("connect database successfully")
}

// GetDBConn 获取DB连接
func GetDBConn() *gorm.DB {
	return conn.DB
}
