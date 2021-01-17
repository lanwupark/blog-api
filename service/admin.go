package service

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/lanwupark/blog-api/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AdminService 管理员服务
type AdminService struct{}

// NewAdminService new
func NewAdminService() *AdminService {
	return &AdminService{}
}

// ErrPageIndexOutOfRange 越界错误
var ErrPageIndexOutOfRange = errors.New("page index out of range")

// UpdateAtFieldName ...
const UpdateAtFieldName = "updateat"

// ArticleQuery 文章查询
func (AdminService) ArticleQuery(query *data.AdminArticleQuery) (*data.ResultListResponse, error) {
	getPageFindOption(&query.PageInfo)
	// 搜索userID
	if query.UserLogin != "" {
		userid, err := userdao.SelectUserIDByUserLogin(query.UserLogin)
		if err != nil {
			return nil, err
		}
		query.UserID = userid
	}
	filter := bson.M{}
	if err := GenerateMongoFilter(filter, query); err != nil {
		return nil, err
	}
	// article
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	articles := []*data.Article{}
	if err := cursor.All(context.TODO(), &articles); err != nil {
		return nil, err
	}
	adminArticleResponses := []*data.AdminArticleResponse{}
	for _, article := range articles {
		resp := new(data.AdminArticleResponse)
		data.DuplicateStructField(article, resp)
		if resp.UserLogin, err = userdao.SelectUserLoginByUserID(resp.UserID); err != nil {
			return nil, err
		}
		adminArticleResponses = append(adminArticleResponses, resp)
	}
	// pageInfo
	query.PageSize = int64(len(adminArticleResponses))
	// total
	if query.Total, err = coll.CountDocuments(context.TODO(), filter); err != nil {
		return nil, err
	}
	return data.NewPageInfoResultListResponse(adminArticleResponses, &query.PageInfo), nil
}

// ArticleUpdate 文章更新
func (AdminService) ArticleUpdate(articleID uint64, status data.CommonType) error {
	return articledao.UpdateArticleStatus(articleID, status)
}

// PhotoQuery 照片查询
// 由于是嵌套结构体查询 这样只有查出所有再排序过滤了...有点耗性能（归纳为经验问题）
func (AdminService) PhotoQuery(query *data.AdminPhotoQuery) (*data.ResultListResponse, error) {
	getPageFindOption(&query.PageInfo)
	// 搜索userID
	if query.UserLogin != "" {
		userid, err := userdao.SelectUserIDByUserLogin(query.UserLogin)
		if err != nil {
			return nil, err
		}
		query.UserID = userid
	}
	filter := bson.M{}
	if err := GenerateMongoFilter(filter, query); err != nil {
		return nil, err
	}
	// photo
	coll := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	albums := []*data.Album{}
	if err := cursor.All(context.TODO(), &albums); err != nil {
		return nil, err
	}
	photoList := []*data.AdminPhotoResponse{}
	for _, album := range albums {
		for _, photo := range album.Photos {
			photoResponse := &data.AdminPhotoResponse{
				AlbumID:           album.AlbumID,
				UserID:            album.UserID,
				AlbumName:         album.Title,
				PhotoName:         photo.Name,
				PhotoOriginalName: photo.OriginalName,
				FileSize:          photo.FileSize,
				Status:            photo.Status,
				CreateAt:          photo.CreateAt,
			}
			if photoResponse.UserLogin, err = userdao.SelectUserLoginByUserID(album.UserID); err != nil {
				return nil, err
			}
			photoList = append(photoList, photoResponse)
		}
	}
	// 排序
	sort.Sort(data.PhotoCreateAtSortDESC(photoList))
	// 计算start:end
	start := (query.PageIndex - 1) * query.PageIndex
	end := start + query.PageSize
	size := int64(len(photoList))
	if end > size {
		end = size
	}
	if start >= size {
		return nil, ErrPageIndexOutOfRange
	}
	res := photoList[start:end]
	// pageInfo
	query.PageSize = int64(len(res))
	// total
	query.Total = size
	return data.NewPageInfoResultListResponse(res, &query.PageInfo), nil
}

// PhotoUpdate 照片更新
func (AdminService) PhotoUpdate(photoName string, status data.CommonType) error {
	return albumdao.UpdatePhotoStatus(photoName, status)
}

// CommentQuery 评论查询
// 由于是嵌套结构体查询 这样只有查出所有再排序过滤了...有点耗性能（归纳为经验问题）
func (AdminService) CommentQuery(query *data.AdminCommentQuery) (*data.ResultListResponse, error) {
	getPageFindOption(&query.PageInfo)
	// 搜索userID
	if query.UserLogin != "" {
		userid, err := userdao.SelectUserIDByUserLogin(query.UserLogin)
		if err != nil {
			return nil, err
		}
		query.UserID = userid
	}
	filter := bson.M{}
	if err := GenerateMongoFilter(filter, query); err != nil {
		return nil, err
	}
	// 嵌套查询
	if query.CommentID != 0 {
		filter["comments.commentid"] = query.CommentID
	}
	// article
	coll := conn.MongoDB.Collection(data.MongoCollectionArticle)
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	articles := []*data.Article{}
	if err := cursor.All(context.TODO(), &articles); err != nil {
		return nil, err
	}
	commemts := []*data.AdminCommentResponse{}
	for _, article := range articles {
		for _, comment := range article.Comments {
			adminCommentResponse := &data.AdminCommentResponse{
				ArticleID: article.ArticleID,
				CommentID: comment.CommentID,
				UserID:    article.UserID,
				Status:    comment.Status,
				Content:   comment.Content,
				CreateAt:  comment.CreateAt,
			}
			if adminCommentResponse.UserLogin, err = userdao.SelectUserLoginByUserID(article.UserID); err != nil {
				return nil, err
			}
			commemts = append(commemts, adminCommentResponse)
		}
	}
	// 排序
	sort.Sort(data.CommentCreateAtSortDESC(commemts))
	// 计算start:end
	start := (query.PageIndex - 1) * query.PageIndex
	end := start + query.PageSize
	size := int64(len(commemts))
	if end > size {
		end = size
	}
	if start >= size {
		return nil, ErrPageIndexOutOfRange
	}
	adminCommentResponses := commemts[start:end]
	// pageInfo
	query.PageSize = int64(len(adminCommentResponses))
	// total
	query.Total = size
	return data.NewPageInfoResultListResponse(adminCommentResponses, &query.PageInfo), nil
}

// CommentUpdate 评论更新
func (AdminService) CommentUpdate(commentID uint64, status data.CommonType) error {
	return articledao.UpdateCommentStatus(commentID, status)
}

// UserQuery 用户查询 整点硬性拼接
func (AdminService) UserQuery(query *data.AdminUserQuery) (*data.ResultListResponse, error) {
	getPageFindOption(&query.PageInfo)
	var sb strings.Builder
	args := []interface{}{}
	sb.WriteString("SELECT * FROM users WHERE 1=1")
	if query.UserID != 0 {
		sb.WriteString(" AND user_id=?")
		args = append(args, query.UserID)
	}
	if query.UserLogin != "" {
		sb.WriteString(" AND user_login=?")
		args = append(args, query.UserLogin)
	}
	if query.Status != "" {
		sb.WriteString(" AND status=?")
		args = append(args, query.Status)
	}
	if query.DateFrom != nil {
		sb.WriteString(" AND update_at>?")
		args = append(args, query.DateFrom)
	}
	if query.DateTo != nil {
		sb.WriteString(" AND update_at<?")
		args = append(args, query.DateTo)
	}
	sb.WriteString(" ORDER BY update_at DESC LIMIT ?,?")
	args = append(args, query.PageIndex-1, query.PageSize)
	db := conn.DB
	users := []*data.User{}
	var total int64
	sql := sb.String()
	countSQL := strings.Replace(sql, "*", "COUNT(0)", 1)
	// count
	if err := db.Get(&total, countSQL, args...); err != nil {
		return nil, err
	}
	// res
	if err := db.Select(&users, sb.String(), args...); err != nil {
		return nil, err
	}
	query.Total = total
	res := []*data.AdminUserResponse{}
	for _, user := range users {
		var userResponse data.AdminUserResponse
		data.DuplicateStructField(user, &userResponse)
		userResponse.Email = user.Email.String
		userResponse.Localtion = user.Location.String
		userResponse.Blog = user.Blog.String
		res = append(res, &userResponse)
	}
	return data.NewPageInfoResultListResponse(res, &query.PageInfo), nil
}

// UserUpdate 用户更新
func (AdminService) UserUpdate(userID uint, status data.CommonType) error {
	return userdao.UpdateUserStatus(userID, status)
}

// GetFeedback 查询反馈
func (AdminService) GetFeedback(pageInfo *data.PageInfo) (*data.ResultListResponse, error) {
	option := getPageFindOption(pageInfo)
	coll := conn.MongoDB.Collection(data.MongoCollectionFeedback)
	cursor, err := coll.Find(context.TODO(), bson.D{}, option)
	if err != nil {
		return nil, err
	}
	feedbacks := []*data.Feedback{}
	if err = cursor.All(context.TODO(), &feedbacks); err != nil {
		return nil, err
	}
	// pageInfo
	pageInfo.PageSize = int64(len(feedbacks))
	// total
	if pageInfo.Total, err = coll.CountDocuments(context.TODO(), bson.D{}); err != nil {
		return nil, err
	}
	return data.NewPageInfoResultListResponse(feedbacks, pageInfo), nil
}

// GenerateMongoFilter 利用反射 src的tag(mongo)的条件来生成filter条件 字段名的小写作为key value作为值
func GenerateMongoFilter(filter bson.M, src interface{}) error {
	// 判断是否是指针类型
	if reflect.TypeOf(src).Kind() != reflect.Ptr {
		return errors.New("the param shoud be a pointer to the struct type")
	}
	// Elem()会去指针指向的值
	if reflect.TypeOf(src).Elem().Kind() != reflect.Struct {
		return errors.New("the param shoud be a pointer to the struct type")
	}
	// 解指针 获取 value
	srcValue := reflect.ValueOf(src).Elem()
	// 解指针 获取 type
	srcType := reflect.TypeOf(src).Elem()

	for i := 0; i < srcType.NumField(); i++ {
		fieldValue := srcValue.Field(i)
		fieldType := srcType.Field(i)
		// 获取匿名结构体
		if fieldType.Anonymous && fieldValue.Kind() == reflect.Struct {
			if err := GenerateMongoFilter(filter, fieldValue.Addr().Interface()); err != nil {
				return err
			}
			continue
		}
		// 变小写
		fieldName := strings.ToLower(fieldType.Name)
		// 根据tag获取条件
		tagName := srcType.Field(i).Tag.Get(data.MongoTag)
		if tagName != "" {
			switch data.MongoCondition(tagName) {
			// equal condition
			case data.MongoEqual:
				if !isDefaultValue(fieldValue) {
					filter[fieldName] = fieldValue.Interface()
				}
			// like condition
			case data.MongoLike:
				if fieldValue.Kind() == reflect.String && fieldValue.String() != "" {
					filter[fieldName] = primitive.Regex{Pattern: fieldValue.String(), Options: "i"}
				}
			// gt condition:
			case data.MongoGreatThan:
				if !isDefaultValue(fieldValue) {
					// 时间结构 特殊逻辑
					if time, ok := fieldValue.Interface().(*time.Time); ok {
						filter[UpdateAtFieldName] = bson.D{{"$gt", time}}
					} else {
						filter[fieldName] = bson.D{{"$gt", fieldValue.Interface()}}
					}

				}
			// lt condition:
			case data.MongoLessThan:
				if !isDefaultValue(fieldValue) {
					if time, ok := fieldValue.Interface().(*time.Time); ok {
						filter[UpdateAtFieldName] = bson.D{{"$lt", time}}
					} else {
						filter[fieldName] = bson.D{{"$lt", fieldValue.Interface()}}
					}
				}
			}

		}

	}
	return nil
}

func isDefaultValue(value reflect.Value) bool {
	kind := value.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint() == 0
	case reflect.String:
		return value.String() == ""
	}
	return value.IsNil()
}

// 设置搜索条件 分页 排序
func getPageFindOption(pageInfo *data.PageInfo) *options.FindOptions {
	if pageInfo.PageIndex < data.DefaultPageIndex {
		pageInfo.PageIndex = data.DefaultPageIndex
	}
	if pageInfo.PageSize <= 0 {
		pageInfo.PageSize = data.DefaultPageSize
	}
	skip := (pageInfo.PageIndex - 1) * pageInfo.PageIndex
	limit := pageInfo.PageSize
	find := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort:  bson.D{{UpdateAtFieldName, -1}}, //逆序
	}
	return find
}
