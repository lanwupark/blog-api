package service_test

import (
	"testing"
	"time"

	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/service"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGenerateFilter(t *testing.T) {
	now := time.Now()
	article := &data.AdminArticleQuery{
		ArticleID: 1111,
		UserLogin: "aaa",
		UserID:    123,
		Title:     "aaa",
		Content:   "aaa",
		DateInterval: data.DateInterval{
			DateFrom: &now,
			DateTo:   &now,
		},
	}
	filter := bson.M{}
	err := service.GenerateMongoFilter(filter, article)
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)
	assert.Equal(uint64(1111), filter["articleid"])
	assert.Equal(uint(123), filter["userid"])
	assert.Nil(filter["status"])
	assert.NotNil(filter["content"])
}
