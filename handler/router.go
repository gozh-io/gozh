package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/handler"
)

/*
	在这里添加gin的路由,需要开发页面的同学添加
*/
func AllRouter(prefix string, router *gin.Engine) {
	//首页
	index := fmt.Sprintf("%s/%s", prefix, "")
	router.GET(index, handler.Demo)

}
