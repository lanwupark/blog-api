package dao

import (
	"context"
	"strconv"

	"github.com/apex/log"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
	"go.mongodb.org/mongo-driver/bson"
)

// AlbumDao ...
type AlbumDao struct{}

// NewAlbumDao ...
func NewAlbumDao() *AlbumDao {
	return &AlbumDao{}
}

// Get 查询相册集
func (AlbumDao) Get(albumID uint64) (*data.Album, error) {
	albumColl := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	result := albumColl.FindOne(context.TODO(), bson.D{{"albumid", albumID}, {"status", data.Normal}})
	if result.Err() != nil {
		return nil, result.Err()
	}
	var album data.Album
	result.Decode(&album)
	return &album, nil
}

// AddPhoto 添加一张照片到集合里
func (AlbumDao) AddPhoto(userID uint, albumID uint64, photo *data.Photo) (int64, error) {
	albumColl := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	res, err := albumColl.UpdateOne(context.TODO(),
		bson.D{{"albumid", albumID}, {"userid", userID}}, //filter
		bson.D{ //update
			{"$push", bson.D{
				{"photos", photo},
			},
			}})
	return res.ModifiedCount, err
}

// CachePhotoInfo 缓存照片信息
func (AlbumDao) CachePhotoInfo(albumID uint64, photo *data.Photo) error {
	rds := conn.Redis
	str, err := util.ToJSONString(photo)
	if err != nil {
		return err
	}
	res := rds.RPush(context.Background(), strconv.Itoa(int(albumID)), str)
	return res.Err()
}

// GetCachePhotoList 获取缓存的相片集合
func (AlbumDao) GetCachePhotoList(albumID uint64) ([]*data.Photo, error) {
	rds := conn.Redis
	// 导出所有相片
	res := rds.LRange(context.TODO(), strconv.Itoa(int(albumID)), 0, -1)
	if res.Err() != nil {
		return nil, res.Err()
	}
	strs, _ := res.Result()
	result := make([]*data.Photo, len(strs))
	for index, val := range strs {
		var photo data.Photo
		// 解析json数据
		if err := util.FromJSONString(val, &photo); err != nil {
			return nil, err
		}
		result[index] = &photo
	}
	return result, nil
}

// AddAlbum 添加相册信息 到mongo里
func (AlbumDao) AddAlbum(album *data.Album) error {
	albumColl := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	_, err := albumColl.InsertOne(context.TODO(), album)
	return err
}

// DelCachePhotoListData 删除缓存的list
func (AlbumDao) DelCachePhotoListData(albumID uint64) error {
	rds := conn.Redis
	res := rds.Del(context.TODO(), strconv.Itoa(int(albumID)))
	return res.Err()
}

// FindByUserID 通过用户id查询
func (AlbumDao) FindByUserID(userID uint) ([]*data.Album, error) {
	albumColl := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	ctx := context.TODO()
	cursor, err := albumColl.Find(ctx, bson.M{"userid": userID, "status": data.Normal})
	if err != nil {
		return nil, err
	}
	res := make([]*data.Album, 0)
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// HitAddition 点击数增加
func (AlbumDao) HitAddition(albumID uint64) {
	coll := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	_, err := coll.UpdateOne(context.TODO(), bson.M{"albumid": albumID}, bson.D{
		{"$inc", bson.D{
			{"hits", 1},
		}},
	})
	if err != nil {
		log.WithError(err).Errorf("add hit error for: %d", albumID)
	}
}

// FindMaintainByUserID 查询maintain user id
func (AlbumDao) FindMaintainByUserID(userID uint) ([]*data.AlbumMaintainResponse, error) {
	albumColl := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	ctx := context.TODO()
	cursor, err := albumColl.Find(ctx, bson.M{"userid": userID, "status": data.Normal})
	if err != nil {
		return nil, err
	}
	res := []*data.AlbumMaintainResponse{}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
