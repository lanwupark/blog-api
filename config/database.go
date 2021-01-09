package config

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//配置mysql,mongodb,redis
var (
	conn = new(Connection)
)

// Connection 数据库连接
type Connection struct {
	DB          *sqlx.DB        // sqlx DB
	mongoClient *mongo.Client   // mongo Client
	MongoDB     *mongo.Database // mongo DB
	Redis       *redis.Client
}

// Config 实现配置接口
func (c *Connection) Config(configs *Configs) {
	// 连mysql
	{
		c.DB = sqlx.MustConnect("mysql", configs.DSN)
		log.Info("Connect to MySQL!")
	}

	// 连mongo
	{
		// Set client options
		clientOptions := options.Client().ApplyURI(configs.MongoURL)
		// 10s 上下文
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			panic(err)
		}
		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			panic(err)
		}
		log.Info("Connected to MongoDB!")
		// 赋值
		c.mongoClient = client
		c.MongoDB = client.Database(configs.MongoDatabase)
	}
	// 连redis
	{
		opt, err := redis.ParseURL(configs.RedisURL)
		if err != nil {
			log.WithError(err).Error("connecting to redis error")
			panic(err)
		}
		rdb := redis.NewClient(opt)
		if status := rdb.Ping(context.TODO()); status.Err() != nil {
			log.WithError(status.Err()).Error("connecting to redis error")
			panic(status.Err())
		}
		c.Redis = rdb
		log.Info("Connected to Redis!")
	}
}

// Shutdown 结束
func (c *Connection) Shutdown() {
	// 关闭 db
	err := c.DB.Close()
	if err != nil {
		log.Errorf("%+v\n", err)
	}
	// 关闭 mongo
	err = c.mongoClient.Disconnect(context.TODO())
	if err != nil {
		log.Errorf("%+v\n", err)
	}
}

// GetConnection 获取连接
func GetConnection() *Connection {
	return conn
}
