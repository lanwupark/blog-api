package main

import (
	"flag"
	"fmt"

	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
)

var (
	c = config.GetConfigs()
)

func init() {
	flag.StringVar(&c.DSN, "dsn", "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local", "dabatase connection string")
	flag.StringVar(&c.BindAdreess, "address", ":8080", "server address")
}

func main() {
	fmt.Println("Hello World!")

	flag.Parse()
	// 配置所有
	c.DoConfigAll()
	conn := config.GetDBConn()
	var user data.User
	conn.First(&user)
	fmt.Println(user)
}
