package dao

import (
	"context"

	"github.com/lanwupark/blog-api/data"
)

type CategoryDao struct{}

func NewCategoryDao() *CategoryDao {
	return &CategoryDao{}
}

func (CategoryDao) InsertOneToMongo(category *data.Category) (id interface{}, err error) {
	mongodb := conn.MongoDB
	coll := mongodb.Collection("categories")
	res, err := coll.InsertOne(context.TODO(), category)
	if err != nil {
		return
	}
	id = res.InsertedID
	return
}
