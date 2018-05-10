package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/module/types"
	"github.com/gozh-io/gozh/module/util"
	"net/http"
)

const (
	AUDIT_USER_FAIL  = -1
	USERNAME_MAX_LEN = 30
	PASSWORD_MAX_LEN = 30
)

//审计提交的user信息是否正确
func auditUser(user *util.User) error {
	username_len := len(user.UserName)
	password_len := len(user.PassWord)
	if username_len == 0 {
		return fmt.Errorf("请确认是否输入用户名,或者请查看是否以post form方式传输的用户名")
	} else if username_len >= USERNAME_MAX_LEN {
		return fmt.Errorf("用户名太长了,最长支持%v个字母", USERNAME_MAX_LEN)
	}
	if password_len == 0 {
		return fmt.Errorf("请确认是否输入密码,或者请查看是否以post form方式传输的密码")
	} else if password_len >= USERNAME_MAX_LEN {
		return fmt.Errorf("密码太长了,最长支持%v个字母", USERNAME_MAX_LEN)
	}
	return nil
}

// 用户注册
// 路由: /user/signup
// Example post a  form (username=manu&password=123)
//curl -X POST http://localhost:8080/user/signup -d "username=admin&password=admin"
func UserSignUp(c *gin.Context) {
	resp := types.Response{}
	user := &util.User{}
	c.Bind(user)
	if err := auditUser(user); err != nil {
		resp.Status, resp.Desc = AUDIT_USER_FAIL, fmt.Sprintf("%v", err)
		c.JSON(http.StatusOK, resp)
		return
	}

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
// 路由: /user/login
// Example post a  form (username=manu&password=123)
//curl -X POST http://localhost:8080/user/login -d "username=admin&password=admin"
func UserLogin(c *gin.Context) {
	resp := types.Response{}
	user := &util.User{}
	c.Bind(user)
	if err := auditUser(user); err != nil {
		resp.Status, resp.Desc = AUDIT_USER_FAIL, fmt.Sprintf("%v", err)
		c.JSON(http.StatusOK, resp)
		return
	}

	mongoUser := util.GetMongoUser()
	err, status := mongoUser.CheckUserPassword(user)
	if status == util.USER_USER_PASSWD_CORRECT {
		resp.Status, resp.Desc = status, fmt.Sprintf("输入的用户%v和密码正确", user.UserName)
	} else if status == util.USER_USER_PASSWD_INCORRECT {
		resp.Status, resp.Desc = status, fmt.Sprintf("输入的用户%v或者密码不正确,请重新输入", user.UserName)
	} else {
		resp.Status, resp.Desc = status, fmt.Sprintf("%v", err)
	}

	session := sessions.Default(c)
	session.Set("userid", user.UserName)
	session.Save()

	c.JSON(http.StatusOK, resp)
}
