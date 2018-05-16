package configure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

//Gin
type Gin struct {
	Mode            string `json:"mode"`
	Host            string `json:"host"`
	Url             string `json:"url"`
	Port            int    `json:"port"`
	Timeout_read_s  int    `json:"timeout_read_s"`
	Timeout_write_s int    `json:"timeout_write_s"`
}

// Log 保存日志配置信息
type Log struct {
	File   string `json:"file"`
	Access string `json:"access"`
}

//es
type Es struct {
	Hosts string `json:"hosts"`
}

//mongo
type Collections_names struct {
	User     string `json:"user"`
	Ariticle string `json:"ariticle"`
}

type Mongo struct {
	Hosts             string            `json:"hosts"`
	Connect_timeout_s int               `json:"connect_timeout_s"`
	Username          string            `json:"username"`
	Passwd            string            `json:"passwd"`
	DatabaseName      string            `json:"database_name"`
	Collections_names Collections_names `json:"collections_names"`
}

//pic_addr
type ImageDatabase struct {
	Hosts             string `json:"hosts"`
	Connect_timeout_s int    `json:"connect_timeout_s"`
	Username          string `json:"username"`
	Passwd            string `json:"passwd"`
	DatabaseName      string `json:"database_name"`
	Collection_prifix string `json:"collection_prifix"`
	Access_url_prefix string `json:"access_url_prefix"`
}

//email
type Email struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

//configure
type configure struct {
	Gin           Gin             `json:"gin"`
	Log           Log             `json:"log"`
	Es            Es              `json:"es"`
	Mongo         Mongo           `json:"mongo"`
	ImageDatabase ImageDatabase   `json:"image_database"`
	WhiteList     map[string]bool `json:"white_list"`
	Email         Email           `json:"email"`
}

var (
	conf      *configure
	conf_once sync.Once
)

//Configure 载入json配置文件
func Configure(ctx context.Context, file string) *configure {
	conf_once.Do(func() {
		conf = &configure{}
		if err := conf.init(file); err != nil {
			log.Fatalln(err)
		}
		conf.readEnv()
	})
	return conf
}

//init 载入json配置文件
func (c *configure) init(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error open file %s fail,%v", file, err)
	}
	defer fd.Close()

	decoder := json.NewDecoder(fd)
	for {
		if err := decoder.Decode(c); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func (c *configure) readEnv() {
	mongo_host := os.Getenv("MONGO_DB_HOST")
	if mongo_host != "" {
		c.Mongo.Hosts = mongo_host
	}
	image_host := os.Getenv("IMAGE_DB_HOST")
	if image_host != "" {
		c.ImageDatabase.Hosts = image_host
	}
}

func (c *configure) String() string {
	js, _ := json.MarshalIndent(c, "", "\t")
	return fmt.Sprintf("%s", js)
}

//得到配置实例
func GetConfigure() *configure {
	return conf
}
