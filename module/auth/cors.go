package auth

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/module/configure"
	"time"
)

//cors用法请参考 https://github.com/gin-contrib/cors
func Cors() gin.HandlerFunc {
	conf := configure.GetConfigure()
	white_list := conf.WhiteList

	// CORS for AllowOriginFunc return true is allowing:
	// - GET and POST methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	f := cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if isAllowed, found := white_list[origin]; found { //在白名单里
				if isAllowed {
					return true //且为true
				}
			}
			return false
		},
		MaxAge: 12 * time.Hour,
	})
	return f
}
