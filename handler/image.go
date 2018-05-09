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
	USERID_NOT_EXIST = -1
)

// 上传单个图片
// 路由: /image/upload
// Example post a  multipart forms
//curl -X POST http://localhost:8080/image/upload -F "file=@./Makefile"  -H "Content-Type: multipart/form-data"
func ImageUpload(c *gin.Context) {
	resp := types.Response{}

	session := sessions.Default(c)
	userid := session.Get("userid")
	username := fmt.Sprintf("%v", userid)
	/*
		if userid == nil {
			resp.Status, resp.Desc = USERID_NOT_EXIST, fmt.Sprintf("%v", "请先登录")
			c.JSON(http.StatusOK, resp)
			return
		}
	*/
	fileHeader, _ := c.FormFile("file")
	image := util.GetMongoImage()
	datas, err := image.Upload(fileHeader, username)

	resp.Status, resp.Desc = types.RET_STATUS_OK, fmt.Sprintf("%v", err)
	resp.Datas = datas

	//其他状态
	c.JSON(http.StatusOK, resp)
}

// 上传多个图片
// 路由: /image/uploads
// Example post a  multipart forms
//curl -X POST http://localhost:8080/image/uploads -F "files=@./Makefile"  -H "Content-Type: multipart/form-data"
func ImageUploads(c *gin.Context) {
	resp := types.Response{}
	session := sessions.Default(c)
	userid := session.Get("userid")
	username := fmt.Sprintf("%v", userid)
	/*
		if userid == nil {
			resp.Status, resp.Desc = USERID_NOT_EXIST, fmt.Sprintf("%v", "请先登录")
			c.JSON(http.StatusOK, resp)
			return
		}
	*/
	form, _ := c.MultipartForm()
	fileHeaders := form.File["files"]
	image := util.GetMongoImage()
	datas, err := image.Uploads(fileHeaders, username)

	resp.Status, resp.Desc = types.RET_STATUS_OK, fmt.Sprintf("%v", err)
	resp.Datas = datas

	c.JSON(http.StatusOK, resp)
}
