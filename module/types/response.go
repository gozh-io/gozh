package types

//接口反向
type (
	Response struct {
		Status int         `json:"status"` //返回状态
		Desc   string      `json:"desc"`   //状态描述
		Datas  interface{} `json:"datas"`  //返回数据
	}
)

const (
	RET_STATUS_OK = 0
)
