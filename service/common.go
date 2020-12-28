package service

import (
	"context"

	"github.com/lanwupark/blog-api/data"
)

// CommonService 小需求写这里
type CommonService struct{}

// NewCommonService ...
func NewCommonService() *CommonService {
	return &CommonService{}
}

// AddFeedback 添加反馈
func (CommonService) AddFeedback(user *data.TokenClaimsSubject, request *data.FeedbackRequest) {
	feedback := &data.Feedback{
		UserID:      user.UserID,
		Description: request.Description,
		Contact:     request.Contact,
		UserLogin:   user.UserLogin,
	}
	coll := conn.MongoDB.Collection(data.MongoCollectionFeedback)
	coll.InsertOne(context.TODO(), &feedback)
}
