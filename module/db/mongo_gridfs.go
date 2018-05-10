package db

import (
	//"crypto/md5"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"io"
	"io/ioutil"
	"mime/multipart"
)

type MongoGridfs struct {
	MongoClient *MongoClient
	Prefix      string //gridfs数据库前缀
	GridFS      *mgo.GridFS
}

//获取一个mongo 的gridfs 连接
func (m *MongoGridfs) Connect() error {
	if err := m.MongoClient.Connect(); err != nil {
		return err
	}
	m.MongoClient.DB()
	fs := "fs"
	if m.Prefix != "" {
		fs = m.Prefix
	}
	m.GridFS = m.MongoClient.Session.Db.GridFS(fs)
	return nil
}

func (m *MongoGridfs) Close() {
	m.MongoClient.Close()
}

//上传每一个返回结果
type RespUpload struct {
	Id       interface{}
	Filename string
}

//上传post 单个文件
func (m *MongoGridfs) Upload(fileHeader *multipart.FileHeader, username string) (*RespUpload, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	//计算md5值
	/*
		hash := md5.New()
		if _, err := io.Copy(hash, src); err != nil {
			return err
		}
		md5_str := fmt.Sprintf("%x", hash.Sum(nil))
	*/
	//整理上传内容
	dst, err := m.GridFS.Create(fileHeader.Filename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	//上传
	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}
	//设置默认metadata,记录每个文件属于哪个用户
	metadata := bson.M{"username": username}
	dst.SetMeta(metadata)

	//获取_id
	respUpload := &RespUpload{}
	respUpload.Id = dst.Id()
	respUpload.Filename = fileHeader.Filename

	return respUpload, nil
}

//上传post 多个文件
func (m *MongoGridfs) Uploads(fileHeaders []*multipart.FileHeader, username string) (map[string]*RespUpload, error) {
	_map := make(map[string]*RespUpload)
	for _, fileHeader := range fileHeaders {
		respUpload, err := m.Upload(fileHeader, username)
		if err != nil {
			return _map, err
		}
		_map[fileHeader.Filename] = respUpload
	}
	return _map, nil
}

//根据id获取单个文件
func (m *MongoGridfs) OpenId(id interface{}) ([]byte, error) {
	id_str := fmt.Sprintf("%v", id)
	objectId := bson.ObjectIdHex(id_str)
	gridfile, err := m.GridFS.OpenId(objectId)
	if err != nil {
		return nil, err
	}
	defer gridfile.Close()
	return ioutil.ReadAll(gridfile)
}

//通过文件名删除,删除文件名是name的所有文件
func (m *MongoGridfs) Remove(name string) error {
	return m.GridFS.Remove(name)
}

//通过id删除,只删除一个文件
func (m *MongoGridfs) RemoveId(id interface{}) error {
	id_str := fmt.Sprintf("%v", id)
	objectId := bson.ObjectIdHex(id_str)
	return m.GridFS.RemoveId(objectId)
}
