package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/module/types"
	"github.com/gozh-io/gozh/module/util"
	"net/http"
)

// 用户注册
func UserSignUp(c *gin.Context) {
	user := &util.User{}
	c.BindJSON(user)
	resp := types.Response{}

	mongoUser := util.GetMongoUser()
	err, status := mongoUser.CheckUser(user)
	if status == util.USER_USER_NOT_EXIST {
		err, status := mongoUser.CreateUser(user)
		if status == util.USER_CREATE_USER_SUCCESS {
			resp.Status, resp.Desc = status, fmt.Sprintf("成功创建用户%v", user.UserName)
		} else if status == util.USER_CREATE_USER_FAIL {
			resp.Status, resp.Desc = status, fmt.Sprintf("创建用户%v失败,请再尝试一下", user.UserName)
		} else {
			resp.Status, resp.Desc = status, fmt.Sprintf("%v", err)
		}
	} else if status == util.USER_USER_EXIST {
		resp.Status, resp.Desc = status, fmt.Sprintf("用户%v已存在,换个用户名试试", user.UserName)
	} else {
		resp.Status, resp.Desc = status, fmt.Sprintf("%v", err)
	}
	//其他状态
	c.JSON(http.StatusOK, resp)
}

// 用户登录
func UserLogin(c *gin.Context) {
	session := sessions.Default(c)
	userid := session.Get("userid")

	user := &util.User{}
	c.BindJSON(user)
	resp := types.Response{}

	mongoUser := util.GetMongoUser()
	err, status := mongoUser.CheckUserPassword(user)
	if status == util.USER_USER_PASSWD_CORRECT {
		resp.Status, resp.Desc = status, fmt.Sprintf("输入的用户%v和密码正确", user.UserName)
	} else if status == util.USER_USER_PASSWD_INCORRECT {
		resp.Status, resp.Desc = status, fmt.Sprintf("输入的用户%v或者密码不正确,请重新输入", user.UserName)
	} else {
		resp.Status, resp.Desc = status, fmt.Sprintf("%v", err)
	}

	session.Set("userid", userid)
	session.Save()

	c.JSON(http.StatusOK, resp)
}
