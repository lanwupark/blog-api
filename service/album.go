package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lanwupark/blog-api/config"
	"github.com/lanwupark/blog-api/data"
	"github.com/lanwupark/blog-api/util"
	"go.mongodb.org/mongo-driver/bson"
)

// AlbumService 相册服务
type AlbumService struct{}

// NewAlbumService ...
func NewAlbumService() *AlbumService {
	return &AlbumService{}
}

// AddPhoto 添加一张照片
// 逻辑 此方法有两种用途 1.写入到已经存储好的相册中 2. 写入到还未创建好的相册当中
// 在保存成功后需要在Mongo里查询 是否有该albumID数据
//    如果有 则生成一条uuid相册信息 插入其中
//    如果没有 则代表新建相册 需要将生成的uuid 存入redis里缓存起来 等待相册创建提交
func (AlbumService) AddPhoto(userID uint, albumID uint64, fileName string, reader io.Reader) (*data.AddPhotoResponse, error) {
	path := config.GetConfigs().FileBaseDir
	maxSize := config.GetConfigs().FileMaxSize
	bufReader := bufio.NewReader(reader)
	// uuid文件名
	uuidFileName := fmt.Sprintf("%s%s", util.NewUUID(), fileName[strings.LastIndex(fileName, "."):])
	// 文件夹+文件
	path = filepath.Join(path, uuidFileName)
	// 检查文件夹是否存在
	dirPath := filepath.Dir(path)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	// 获取文件信息 存在着删除它
	_, err = os.Stat(path)
	// 逻辑:如果错误为空则代表能检测到该文件
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			return nil, fmt.Errorf("Unable to delete file:%v", err)
		}
	} else if !os.IsNotExist(err) {
		// 如果是除了文件不存在的错误
		return nil, fmt.Errorf("Unable to get file:%v", err)
	}

	// 创建新文件
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to create file: %v", err)
	}
	defer f.Close()
	// 读写文件
	bufWriter := bufio.NewWriter(f)
	size := 0
	buf := make([]byte, 4096)
	// 由于没有找到一个可以事先预知文件大小的方法 我就采用的边读边写的方法来计算文件大小 如果文件大小超过限制 那么会关闭文件并删除它
	for {
		// 如果超出限制
		if size > maxSize {
			f.Close()
			os.Remove(path)
			return nil, fmt.Errorf("file size must be less than %d byte", maxSize)
		}
		// read
		n, err := bufReader.Read(buf)
		// 文件读取到末尾
		if err == io.EOF {
			break
		}
		// 其它错误
		if err != nil {
			return nil, fmt.Errorf("Unable to read file: %v", err)
		}
		// write
		_, err = bufWriter.Write(buf[:n])
		if err != nil {
			return nil, fmt.Errorf("Unable to write file: %v", err)
		}
		size += n
	}
	// 用户上传的空文件 虚晃一枪
	if size == 0 {
		f.Close()
		os.Remove(path)
		return nil, errors.New("empty request body")
	}
	// 尝试插入到相册集里的photo数据
	photo := &data.Photo{
		Name:         uuidFileName,
		OriginalName: fileName,
		FileSize:     uint64(size),
		Status:       data.Normal,
		CreateAt:     time.Now(),
		UpdateAt:     time.Now(),
	}
	if count, err := albumdao.AddPhoto(userID, albumID, photo); err != nil {
		return nil, err
	} else if count == 0 {
		// 没有该集合 先放入redis里缓存起来
		if err = albumdao.CachePhotoInfo(albumID, photo); err != nil {
			return nil, err
		}
	}
	// 回写数据
	resp := &data.AddPhotoResponse{
		FileName:         uuidFileName,
		OriginalFileName: fileName,
		FileSize:         int(size),
	}
	// 终于上传成功了
	return resp, nil
}

// NewAlbum 新建相册
// 逻辑:
// 1. 从redis里读出该albumID的相册集合
// 2. 获取用户真正要上传的图片uuid 删除不需要的图片
// 3. 添加数据到mongo
func (AlbumService) NewAlbum(userID uint, albumReq *data.AddAlbumRequest) error {
	dir := config.GetConfigs().FileBaseDir
	list, err := albumdao.GetCachePhotoList(albumReq.AlbumID)
	if err != nil {
		return err
	}
	// 删除缓存的数据
	if err = albumdao.DelCachePhotoListData(albumReq.AlbumID); err != nil {
		return err
	}
	photoSet := make(map[string]bool, len(albumReq.PhotoList))
	// 剩下的 photo list
	realPhotoList := make([]*data.Photo, 0)
	for _, val := range albumReq.PhotoList {
		path := filepath.Join(dir, val)
		_, err := os.Stat(path)
		// 检测请求里的数据集合文件是否真正存在 若不存在则忽略该文件
		if os.IsNotExist(err) {
			continue
		}
		photoSet[val] = true
	}
	for _, val := range list {
		if _, ok := photoSet[val.Name]; ok {
			realPhotoList = append(realPhotoList, val)
		} else {
			// 删除多余的照片
			path := filepath.Join(dir, val.Name)
			if err = os.Remove(path); err != nil {
				return err
			}
		}
	}
	// 设置封面 如果covername未设置 并且如果集合不为空 那么取第一个
	if albumReq.CoverName == "" && len(realPhotoList) > 0 {
		albumReq.CoverName = realPhotoList[0].Name
	}
	album := &data.Album{
		AlbumID:     albumReq.AlbumID,
		UserID:      userID,
		Title:       albumReq.Title,
		Description: albumReq.Description,
		CoverName:   albumReq.CoverName,
		Location:    albumReq.Location,
		Hits:        0,
		Status:      data.Normal,
		Photos:      realPhotoList,
		CreateAt:    time.Now(),
		UpdateAt:    time.Now(),
	}
	err = albumdao.AddAlbum(album)
	return err
}

// CancelNewAlbum 取消新建相册
func (AlbumService) CancelNewAlbum(albumID uint64) error {
	// 删除照片
	if list, err := albumdao.GetCachePhotoList(albumID); err == nil {
		for _, val := range list {
			path := filepath.Join(config.GetConfigs().FileBaseDir, val.Name)
			if err = os.Remove(path); err != nil {
				return err
			}
		}
	} else {
		return err
	}
	// 删除缓存集合
	return albumdao.DelCachePhotoListData(albumID)
}

// EditAlbum 编辑相册
func (AlbumService) EditAlbum(userID uint, albumID uint64, albumReq *data.EditAlbumRequest) error {
	filter := bson.D{{"userid", userID}, {"albumid", albumID}}
	update := bson.M{"updateat": time.Now()}
	if strings.TrimSpace(albumReq.Title) != "" {
		update["title"] = albumReq.Title
	}
	if strings.TrimSpace(albumReq.Description) != "" {
		update["description"] = albumReq.Description
	}
	if strings.TrimSpace(albumReq.Location) != "" {
		update["localtion"] = albumReq.Location
	}
	if strings.TrimSpace(albumReq.CoverName) != "" {
		update["covername"] = albumReq.CoverName
	}
	albumColl := conn.MongoDB.Collection(data.MongoCollectionAlbum)
	ctx := context.TODO()
	// 更新数据
	if _, err := albumColl.UpdateOne(ctx, filter, bson.D{{"$set", update}}); err != nil {
		return err
	}
	// 删除相片
	for _, val := range albumReq.DeletePhotoList {
		path := filepath.Join(config.GetConfigs().FileBaseDir, val)
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	// 删除mongo里数据
	pull := bson.D{
		{"$pull", bson.D{
			{"photos", bson.D{
				{"name", bson.D{
					{"$in", albumReq.DeletePhotoList},
				}},
			},
			}},
		},
	}
	_, err := albumColl.UpdateOne(ctx, bson.D{{"albumid", albumID}, {"userid", userID}}, pull)
	return err
}

// GetAlbumInfo 获取相册信息
func (AlbumService) GetAlbumInfo(userID uint) ([]*data.AlbumMaintainResponse, error) {
	albums, err := albumdao.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	albumMaintains := []*data.AlbumMaintainResponse{}
	for _, val := range albums {
		albumMaintain := &data.AlbumMaintainResponse{
			AlbumID:   val.AlbumID,
			AlbumName: val.Title,
			CoverName: val.CoverName,
			Location:  val.Location,
			CreateAt:  val.CreateAt,
		}
		albumMaintains = append(albumMaintains, albumMaintain)
	}
	return albumMaintains, nil
}

// GetAlbumDetail 获取相册详细信息
func (AlbumService) GetAlbumDetail(albumID uint64) (*data.Album, error) {
	resp, err := albumdao.Get(albumID)
	if err != nil {
		return nil, err
	}
	// 过滤被拉黑的相片
	var photos []*data.Photo
	for _, val := range resp.Photos {
		if val.Status == data.Normal {
			photos = append(photos, val)
		}
	}
	resp.Photos = photos
	// 增加点击
	go albumdao.HitAddition(albumID)
	return resp, nil
}
