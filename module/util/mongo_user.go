package util

import (
	"github.com/gozh-io/gozh/module/db"
	"sync"
	"github.com/gozh-io/gozh/module/configure"
	"fmt"
	"github.com/gozh-io/gozh/module/data"
	"github.com/globalsign/mgo/bson"
	"log"
)

type mongoUser struct {
	Mongo *db.MongoClient
}

var (
	user      *mongoUser
	user_once sync.Once
)

//对外方法,使用时,先init,再start,退出时stop
func (m *mongoUser) Init() error {
	conf := configure.GetConfigure()
	mongo := &db.MongoClient{
		Hosts:                conf.Mongo.Hosts,
		Database:             conf.Mongo.DatabaseName,
		Collection:           "user",
		ConnectTimeoutSecond: conf.Mongo.Connect_timeout_s,
	}
	m.Mongo = mongo
	return nil
}

func (m *mongoUser) NewUser(user, password string) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	new_user := data.User{ID: GetUserId(user), User: user, Password: cryptoSha256(password)}

	if err := collection.Insert(new_user); err != nil {
		return fmt.Errorf("create new user %v %v", user, err)
	}
	return nil
}

func (m *mongoUser) CheckUserExist(user string) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	// 查询user字段是否存在
	q := bson.M{"user": user}

	count, err := collection.Find(q).Count()
	// 查询出错处理
	if err != nil {
		return fmt.Errorf("chech %v fail, %v", user, err)
	}
	// user存在处理
	if count != 0 {
		return fmt.Errorf("user %v is exist", user)
	}
	return nil
}

func (m *mongoUser) CheckUserAndPassword(user, password string) (int, error) {
	if err := m.Mongo.Connect(); err != nil {
		return -1, fmt.Errorf("connect db fail, %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	q := bson.M{"user": user, "password": cryptoSha256(password)}
	count, err := collection.Find(q).Count()

	if err != nil {
		return -1, fmt.Errorf("find %v fail, %v", user, err)
	}
	return count, nil
}

// 根据用户名获取用户的id
func (m *mongoUser) GetUserID(username string) (id string, err error) {
	if err = m.Mongo.Connect(); err != nil {
		id = ""
		err = fmt.Errorf("connect db fail, %v", err)
		return
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	q := bson.M{"user": username}
	var users data.Users
	collection.Find(q).All(&users)
	if len(users) > 1 {
		id = ""
		err = fmt.Errorf("the number of user is too many, count: %v", len(users))
		return
	}
	id = users[0].ID.String()
	err = nil
	return
}

func GetUser() *mongoUser {
	user_once.Do(func() {
		user = &mongoUser{}
		if err := user.Init(); err != nil {
			log.Fatalln(err)
		}
	})
	return user
}


