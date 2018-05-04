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

	/*
	  这里添加其他 url和handler处理关系
	*/
}
