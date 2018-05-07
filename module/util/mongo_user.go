package util

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/gozh-io/gozh/module/common"
	"github.com/gozh-io/gozh/module/configure"
	"github.com/gozh-io/gozh/module/db"
	"log"
	"sync"
)

//状态
const (
	USER_USER_EXIST     = -1
	USER_USER_NOT_EXIST = -2
	//
	USER_USER_PASSWD_CORRECT   = 0 //登陆成功后返回状态
	USER_USER_PASSWD_INCORRECT = -3
	//
	USER_CREATE_USER_SUCCESS = 0 //注册成功后返回状态
	USER_CREATE_USER_FAIL    = -4
	//
	USER_CONNECT_MONGO_FAIL = -5
	USER_FIND_USER_FAIL     = -6
)

type User struct {
	UserName   string `form:"username" json:"username"`
	PassWord   string `form:"password" json:"password"`
	CreateTime string `json:"createtime"`
}

type mongoUser struct {
	Mongo *db.MongoClient
	ctx   context.Context
}

var (
	user      *mongoUser
	user_once sync.Once
)

//对外方法,使用时,先init,,再start,退出时stop
func (m *mongoUser) init() error {
	conf := configure.GetConfigure()
	mongo := &db.MongoClient{
		Hosts:                conf.Mongo.Hosts,
		Database:             conf.Mongo.DatabaseName,
		Collection:           conf.Mongo.Collections_names.User,
		ConnectTimeoutSecond: conf.Mongo.Connect_timeout_s,
	}
	m.Mongo = mongo
	return nil
}

//新建一个用户
//传入user,user结构字段增加时,具有可扩展性
func (m *mongoUser) CreateUser(user *User) (err error, status int) {
	if err := m.Mongo.Connect(); err != nil {
		return err, USER_CONNECT_MONGO_FAIL
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	user.PassWord = common.Sha256Encode(user.PassWord)
	user.CreateTime = common.GetBeiJingDT()

	collection := m.Mongo.GetCollection()
	if err := collection.Insert(user); err != nil {
		return err, USER_CREATE_USER_FAIL
	}
	return nil, USER_CREATE_USER_SUCCESS
}

//检测用户是否存在
//传入user,user结构字段增加时,具有可扩展性
func (m *mongoUser) CheckUser(user *User) (err error, status int) {
	if err := m.Mongo.Connect(); err != nil {
		return err, USER_CONNECT_MONGO_FAIL
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	// 查询username字段是否存在
	userName := user.UserName
	q := bson.M{"username": userName}

	collection := m.Mongo.GetCollection()
	count, err := collection.Find(q).Count()
	// 查询出错处理
	if err != nil {
		return err, USER_FIND_USER_FAIL
	}
	// username存在处理
	if count > 0 {
		return nil, USER_USER_EXIST
	}
	return nil, USER_USER_NOT_EXIST
}

//在用户存在的前提下, 检测用户和密码是否正确
//传入user,user结构字段增加时,具有可扩展性
//当err=nil 且 status=成功
func (m *mongoUser) CheckUserPassword(user *User) (err error, status int) {
	if err := m.Mongo.Connect(); err != nil {
		return err, USER_CONNECT_MONGO_FAIL
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	user.PassWord = common.Sha256Encode(user.PassWord)
	q := bson.M{"username": user.UserName, "password": user.PassWord}

	collection := m.Mongo.GetCollection()
	count, err := collection.Find(q).Count()
	if err != nil {
		return err, USER_FIND_USER_FAIL
	}
	if count == 0 { //没找到
		return nil, USER_USER_PASSWD_INCORRECT
	}
	return nil, USER_USER_PASSWD_CORRECT
}

//对外使用方法,在main.go文件中,Init中调用
//这里context要对单实例进行控制,这里mongoUser比较简单,暂时用不到context控制逻辑
func MongoUser(ctx context.Context) *mongoUser {
	user_once.Do(func() {
		user = &mongoUser{}
		if err := user.init(); err != nil {
			log.Fatalln(err)
		}
		user.ctx = ctx
	})
	return user
}

//在其他地方都直接调用,初始化一次,到处使用
func GetMongoUser() *mongoUser {
	return user
}
