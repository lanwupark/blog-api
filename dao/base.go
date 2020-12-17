package dao

import "github.com/lanwupark/blog-api/config"

var (
	conn = config.GetConnection()
)
