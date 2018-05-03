package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gozh-io/gozh/handler"
	"github.com/gozh-io/gozh/module"
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
	//
	log.SetFlags(log.LstdFlags)
	log.Println("[Main] Starting program")
	defer log.Println("[Main] Exit program successful.")
	//
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //cancel when exit
	//
	Init(ctx, os.Args[1])

	//
	conf := module.GetConfigure()
	g := conf.Gin
	mode := g.Mode
	host := g.Host
	url := g.Url
	port := g.Port
	timeout_read := g.Timeout_read_s
	timeout_write := g.Timeout_write_s

	//配置gin
	gin.SetMode(mode)
	router := gin.New()
	router.Use(module.Logger())
	//添加session管理
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	prefix := fmt.Sprintf("%s", url)
	handler.AllRouter(prefix, router)

	//起一个http服务器
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

func Init(ctx context.Context, filename string) {
	SetupCPU()
	module.Configure(ctx, filename)
	module.Mylog(ctx)
}

func SetupCPU() {
	num := runtime.NumCPU()
	runtime.GOMAXPROCS(num)
}
