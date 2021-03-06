package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/*
	在这里添加gin的路由,需要开发页面的同学添加
*/
func AllRouter(router *gin.Engine, prefix string) {
	//首页
	index := fmt.Sprintf("%s/%s", prefix, "")
	router.GET(index, Demo)

	//user
	user := fmt.Sprintf("%s/%s", prefix, "user")
	userGroup := router.Group(user)
	{
		sign_up := "/signup"
		userGroup.POST(sign_up, UserSignUp)
		login := "/login"
		userGroup.POST(login, UserLogin)
	}
	//image
	image := fmt.Sprintf("%s/%s", prefix, "image")
	imageGroup := router.Group(image)
	{
		imageGroup.POST("/upload", ImageUpload)
		imageGroup.POST("/uploads", ImageUploads)
		imageGroup.GET("/get/:imgid", ImageOpenId)
		imageGroup.GET("/remove/:imgid", ImageRemoveId)
	}
	//captcha
	captcha := fmt.Sprintf("%s/%s", prefix, "captcha")
	router.GET(captcha, GenerateCaptchaHandler)
	//email captcha
	emailCaptcha := fmt.Sprintf("%s/%s", prefix, "email_captcha")
	router.POST(emailCaptcha, EmailCaptcha)
	/*
	  这里添加其他 url和handler处理关系, 请保持这个在最下面
	*/
}
