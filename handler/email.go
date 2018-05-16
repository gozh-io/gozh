package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/module/common"
	"github.com/gozh-io/gozh/module/types"
	"github.com/gozh-io/gozh/module/util"
	"net/http"
)

// status
const (
	// send email status
	SEND_EMAIL_OK    = 0
	SEND_EMAIL_FAILD = -1
	// verfiy email captcha status
	EMAIL_CAPTCHA_IS_RIGHT = 0
	EMAIL_CAPTCHA_IS_ERROR = -2
)

// 发送邮箱验证码
// 路由: /email_captcha
// Example post a  form (email=gozh@gozh.io)
//curl -X POST http://localhost:8080/email_captcha -d "email=gozh@gozh.io"
func EmailCaptcha(c *gin.Context) {
	email := c.PostForm("email")
	session := sessions.Default(c)
	emailCaptcha := common.GenerateRandomString(6)
	resp := types.Response{}

	go util.SendEmailCaptchaTo(email, emailCaptcha)

	session.Set("email_captcha", emailCaptcha)
	session.Save()

	resp.Status, resp.Desc = SEND_EMAIL_OK, "邮箱验证码发送成功"
	c.JSON(http.StatusOK, resp)
}

// 验证 邮箱验证码是否正确
// emailcaptcha: 存于session中的email_captcha
// verifyValue: 客户端发来的验证码
func VerfiyEmailCaptcha(emailCaptcha, verifyValue string) error {
	fmt.Println("session",emailCaptcha, "post", verifyValue)
	if emailCaptcha != verifyValue {
		return fmt.Errorf("邮箱验证码错误")
	}
	return nil
}
