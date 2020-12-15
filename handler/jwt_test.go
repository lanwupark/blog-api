package handler_test

import (
	"testing"

	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/handler"
)

func TestJWTSuccess(t *testing.T) {
	user := &data.User{
		UserID:    123456,
		UserLogin: "eanson023",
		IsAdmin:   true,
	}
	token, err := handler.CreateToken(user)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token:%s\n", token)
	parseUser, err := handler.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("parse user:%v\n", parseUser)
}
