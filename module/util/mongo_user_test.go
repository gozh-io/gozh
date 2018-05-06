package util

import (
	"testing"
	"context"
	"github.com/gozh-io/gozh/module/configure"
)

func TestGetUser(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	configure.Configure(ctx, "../../conf/cf.json")
	user := GetUser()
	if user == nil {
		t.Error("TestGetUser is faild")
	} else {
		t.Log("TestGetUser is pass")
	}
}

func TestMongoUser_NewUser(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	configure.Configure(ctx, "../../conf/cf.json")
	user := GetUser()
	if user == nil {
		t.Error("TestMongoUser_NewUser is faild --- user is nil")
	}
	if err := user.NewUser("gozh", "123456"); err != nil {
		t.Error("TestMongoUser_NewUser is faild ---- ", err)
	} else {
		t.Log("TestMongoUser_NewUser is pass")
	}
}

func TestMongoUser_CheckUserExist(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	configure.Configure(ctx, "../../conf/cf.json")
	user := GetUser()
	if user == nil {
		t.Error("TestMongoUser_CheckUserExist is faild --- user is nil")
	}
	if err := user.CheckUserExist("gozh"); err == nil {
		t.Error("TestMongoUser_CheckUserExist is faild ---- ", err)
	} else {
		t.Log("TestMongoUser_CheckUserExist is pass")
	}
}

func TestMongoUser_CheckUserAndPassword(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	configure.Configure(ctx, "../../conf/cf.json")
	user := GetUser()
	if user == nil {
		t.Error("TestMongoUser_CheckUserAndPassword is faild --- user is nil")
	}
	if _,err := user.CheckUserAndPassword("gozh", "123456"); err == nil {
		t.Error("TestMongoUser_CheckUserAndPassword is faild ---- ", err)
	} else {
		t.Log("TestMongoUser_CheckUserAndPassword is pass")
	}
}