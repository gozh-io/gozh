package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/handler"
	"github.com/gozh-io/gozh/module/configure"
	"github.com/gozh-io/gozh/module/mylog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func Usage(program string) {
	fmt.Printf("\nusage: %s conf/cf.json\n", program)
	fmt.Printf("\nconf/cf.json      configure file\n")
}

func main() {
	if len(os.Args) != 2 {
		Usage(os.Args[0])
		os.Exit(-1)
	}
	//设置官方日志包log输出格式
	log.SetFlags(log.LstdFlags)
	log.Println("[Main] Starting program")
	defer log.Println("[Main] Exit program successful.")
	//创建一个context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //send cancel operator when exit
	//在启动http服务前,加载单实例资源
	Init(ctx, os.Args[1])

	//配置gin
	conf := configure.GetConfigure()
	g := conf.Gin
	mode, host, url, port, timeout_read, timeout_write := g.Mode, g.Host, g.Url, g.Port, g.Timeout_read_s, g.Timeout_write_s

	gin.SetMode(mode)
	router := gin.New()
	useMiddleware(router)                     //配置使用中间件
	allRouter(router, fmt.Sprintf("%s", url)) //配置路由

	//起一个goroutine 跑http服务
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      router,
		ReadTimeout:  time.Duration(timeout_read) * time.Second,
		WriteTimeout: time.Duration(timeout_write) * time.Second,
	}
	go func(s *http.Server) {
		log.Printf("[Main] http server start\n")
		err := s.ListenAndServe()
		log.Printf("[Main] http server stop (%+v)\n", err)
	}(s)
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	for {
		select {
		case sig := <-signals:
			log.Println("[Main] Catch signal", sig)
			//平滑关闭server
			err := s.Shutdown(context.Background())
			log.Printf("[Main] start gracefully shuts down http serve %+v", err)
			return
		}
	}
}

//配置gin 的路由
func allRouter(router *gin.Engine, prefix string) {
	handler.AllRouter(router, prefix) //真正的url和handler对应关系在 handler/router.go里配置
}

//配置gin 使用哪些中间件
func useMiddleware(router *gin.Engine) {
	//输出访问日志
	router.Use(mylog.Logger())
	//添加session管理
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	/*
	  这里添加其他中间件
	*/
}

//这里配置程序启动时,需要加载的单实例模块
func Init(ctx context.Context, filename string) {
	SetupCPU()
	configure.Configure(ctx, filename)
	mylog.Mylog(ctx)

	/*
	  这里添加其他单实例
	*/
}

//配置程序使用几个cpu
func SetupCPU() {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)
}
