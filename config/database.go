package config

import (
	"github.com/apex/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	conn = new(Connection)
)

func init() {
	c.RegisterService(conn)
}

// Connection 数据库连接
type Connection struct {
	DB *gorm.DB
}

// Config 实现配置接口
func (c *Connection) Config(configs *Configs) {
	db, err := gorm.Open(mysql.Open(configs.DSN), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	log.Info("connect database successfully")
	c.DB = db
}

// Shutdown 结束
func (c *Connection) Shutdown() {
	sqlDB, err := conn.DB.DB()
	if err != nil {
		log.WithError(err).Error("")
	}
	sqlDB.Close()
}

// GetConnection 获取DB连接
func GetConnection() *Connection {
	return conn
}
