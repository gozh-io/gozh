package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/module/util"
	"net/http"
	"github.com/gin-contrib/sessions"
	"time"
	"log"
)

// 用于handler使用的struct
// 便于gin.Context 的BindJson方法使用

type Text struct {
	Text string `json:"text"`
}

type UserAndPassword struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// 用户注册
func UserSignUp(c *gin.Context) {
	user := &UserAndPassword{}
	c.BindJSON(user)
	mongoUser := util.GetUser()
	err := mongoUser.CheckUserExist(user.User)
	if err != nil {
		c.JSON(http.StatusOK, &Text{Text: "用户名存在"})
		return
	}
	err = mongoUser.NewUser(user.User, user.Password)
	if err != nil {
		log.Printf("[handler.user.UserSignUp] %v : %v \n", time.Now(), err)
		c.JSON(http.StatusOK, &Text{Text: "创建用户名失败"})
		return
	} else {
		log.Printf("[handler.user.UserSignUp] %v create user %s successful \n", time.Now(), user.User)
		c.JSON(http.StatusCreated, &Text{Text: "ok"})
		return
	}
}

// 用户登录
func UserLogin(c *gin.Context) {
	session := sessions.Default(c)
	user_id_session := session.Get("user_id")
	if user_id_session != nil {
		c.JSON(http.StatusOK, &Text{Text: "用户已登录"})
		return
	}
	user := &UserAndPassword{}
	c.BindJSON(user)
	mongoUser := util.GetUser()
	count, err := mongoUser.CheckUserAndPassword(user.User, user.Password)
	if err != nil {
		c.Error(err)
		return
	}
	if count != 1 {
		c.JSON(http.StatusOK, &Text{Text: "用户名或密码错误"})
		return
	}
	user_id, err := mongoUser.GetUserID(user.User)
	if err != nil {
		c.Error(err)
		return
	}
	session.Set("user_id", user_id)
	session.Save()
	log.Printf("user %s login", user.User)
	c.JSON(http.StatusOK, &Text{Text: "ok"})
}
