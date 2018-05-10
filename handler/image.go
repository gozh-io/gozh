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

func checkSession(c *gin.Context) string {
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
	return username
}

// 上传单个图片
// 路由: /image/upload
// Example post a  multipart forms
//curl -X POST http://localhost:8080/image/upload -F "file=@./Makefile"  -H "Content-Type: multipart/form-data"
func ImageUpload(c *gin.Context) {
	username := checkSession(c)
	fileHeader, _ := c.FormFile("file")
	image := util.GetMongoImage()
	datas, err := image.Upload(fileHeader, username)

	resp := types.Response{}
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
	username := checkSession(c)
	form, _ := c.MultipartForm()
	fileHeaders := form.File["files"]
	image := util.GetMongoImage()
	datas, err := image.Uploads(fileHeaders, username)

	resp := types.Response{}
	resp.Status, resp.Desc = types.RET_STATUS_OK, fmt.Sprintf("%v", err)
	resp.Datas = datas

	c.JSON(http.StatusOK, resp)
}

// 获取图片
// 路由: /image/get/:imgid
//curl -X GET http://localhost:8080/image/get/5af2ac59ea83be000581fdde, 这里5af2ac59ea83be000581fdde是img在mongo中的id
func ImageOpenId(c *gin.Context) {
	username := checkSession(c)
	imgid := c.Param("imgid")

	image := util.GetMongoImage()
	datas, err := image.OpenId(imgid, username)
	if err != nil {
		resp := types.Response{}
		resp.Status, resp.Desc = types.RET_STATUS_OK, fmt.Sprintf("%v", err)
		c.JSON(http.StatusOK, resp)
		return
	}
	c.JSON(http.StatusOK, datas)
}

// 删除图片
// 路由: /image/remove/:imgid
//curl -X GET http://localhost:8080/image/remove/5af2ac59ea83be000581fdde, 这里5af2ac59ea83be000581fdde是img在mongo中的id
func ImageRemoveId(c *gin.Context) {
	username := checkSession(c)
	imgid := c.Param("imgid")

	image := util.GetMongoImage()
	err := image.RemoveId(imgid, username)

	resp := types.Response{}
	resp.Status = types.RET_STATUS_OK
	if err != nil {
		resp.Desc = fmt.Sprintf("%v", err)
	} else {
		resp.Desc = "成功删除"
	}
	c.JSON(http.StatusOK, resp)
}
