package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/module/util"
	"github.com/mojocn/base64Captcha"
)

const (
	CAPTCHA_IS_RIGHT = 0
	CAPTCHA_IS_ERROR = -7
)

// 获取验证码图片
// 路由: /captcha
func GenerateCaptchaHandler(c *gin.Context) {
	// get session
	session := sessions.Default(c)
	captchaConfig := util.GetCaptchaConfig()
	//create base64 encoding captcha
	//创建base64图像验证码
	config := captchaConfig.ConfigCharacter
	//GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uiid.
	captchaId, digitCap := base64Captcha.GenerateCaptcha(captchaConfig.Id, config)
	base64Png := base64Captcha.CaptchaWriteToBase64Encoding(digitCap)
	session.Set("captchaId", captchaId)
	session.Save()
	c.String(http.StatusOK, base64Png)
}

//  验证 验证码是否正确
// captchaId: 存于session中
// verifyValue: 客户端发来的验证码
func VerfiyCaptcha(captchaId, verifyValue string) (int, error) {
	verifyResult := base64Captcha.VerifyCaptcha(captchaId, verifyValue)
	if verifyResult {
		return CAPTCHA_IS_RIGHT, nil
	} else {
		return CAPTCHA_IS_ERROR, fmt.Errorf("captcha is error")
	}
}
