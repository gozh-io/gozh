package util

import (
	"context"
	"fmt"
	"github.com/gozh-io/gozh/module/configure"
	"github.com/gozh-io/gozh/module/db"
	"log"
	"mime/multipart"
	"sync"
)

type mongoImage struct {
	Mongo           *db.MongoGridfs
	ctx             context.Context
	AccessUrlPrefix string
}

var (
	image      *mongoImage
	image_once sync.Once
)

//对外方法,使用时,先init,,再start,退出时stop
func (m *mongoImage) init(ctx context.Context) error {
	conf := configure.GetConfigure()
	mongo := &db.MongoGridfs{
		MongoClient: &db.MongoClient{
			Hosts:                conf.ImageDatabase.Hosts,
			Database:             conf.ImageDatabase.DatabaseName,
			ConnectTimeoutSecond: conf.ImageDatabase.Connect_timeout_s,
		},
		Prefix: conf.ImageDatabase.Collection_prifix,
	}
	m.Mongo = mongo
	m.AccessUrlPrefix = conf.ImageDatabase.Access_url_prefix
	m.ctx = ctx
	return nil
}

//上传post 单个文件
func (m *mongoImage) Upload(fileHeader *multipart.FileHeader, username string) (*db.RespUpload, error) {
	if err := m.Mongo.Connect(); err != nil {
		return nil, err
	}
	defer m.Mongo.Close()
	return m.Mongo.Upload(fileHeader, username)
}

//上传post 多个文件
func (m *mongoImage) Uploads(fileHeaders []*multipart.FileHeader, username string) (map[string]*db.RespUpload, error) {
	if err := m.Mongo.Connect(); err != nil {
		return nil, err
	}
	defer m.Mongo.Close()
	return m.Mongo.Uploads(fileHeaders, username)
}

//根据id读文件
func (m *mongoImage) OpenId(id interface{}, username string) ([]byte, error) {
	if err := m.Mongo.Connect(); err != nil {
		return nil, err
	}
	defer m.Mongo.Close()
	return m.Mongo.OpenId(id)
}

//根据id删除文件
func (m *mongoImage) RemoveId(id interface{}, username string) error {
	if err := m.Mongo.Connect(); err != nil {
		return err
	}
	defer m.Mongo.Close()
	return m.Mongo.RemoveId(id)
}

//获取到图片的最终访问地址
func (m *mongoImage) GetAccessUrl(idOrFilename interface{}) string {
	return fmt.Sprintf("%v/%v", m.AccessUrlPrefix, idOrFilename)
}

//对外使用方法,在main.go文件中,Init中调用
//这里context要对单实例进行控制,这里mongoImage比较简单,暂时用不到context控制逻辑
func MongoImage(ctx context.Context) *mongoImage {
	image_once.Do(func() {
		image = &mongoImage{}
		if err := image.init(ctx); err != nil {
			log.Fatalln(err)
		}
	})
	return image
}

//在其他地方都直接调用,初始化一次,到处使用
func GetMongoImage() *mongoImage {
	return image
}
