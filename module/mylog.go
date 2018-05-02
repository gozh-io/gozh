package module

import (
	"context"
	"github.com/cihub/seelog"
	"log"
	"sync"
)

type mylog struct {
	File   string //普通日志
	File_a string //gin访问日志

	Log   seelog.LoggerInterface
	Log_a seelog.LoggerInterface

	StructName string
}

var (
	my_log     *mylog
	mylog_once sync.Once
)

// Mylog  创建mylog单实例
func Mylog(ctx context.Context) *mylog {
	mylog_once.Do(func() {
		my_log = &mylog{}
		if err := my_log.loadConfigure(); err != nil {
			log.Fatalln(err)
		}
		my_log.waitStop(ctx)
	})
	return my_log
}

// LoadConfigure 从file里读取seelog 配置
func (l *mylog) loadConfigure() error {
	conf := GetConfigure()
	file := conf.Log.File
	file_a := conf.Log.Access

	logger, err := seelog.LoggerFromConfigAsFile(file)
	if err != nil {
		return err
	}
	l.Log = logger

	logger_a, err := seelog.LoggerFromConfigAsFile(file_a)
	if err != nil {
		return err
	}
	l.Log_a = logger_a

	l.StructName = GetStructName(l)
	log.Printf("[%s] Start\n", l.StructName)
	return nil
}

func (l *mylog) waitStop(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				l.Log_a.Flush()
				l.Log.Flush()
				log.Printf("[%s] Stoped\n", l.StructName)
				return
			}
		}
	}()
}

// Infof 输出info信息
func (l *mylog) Access() seelog.LoggerInterface {
	return l.Log_a
}

func (l *mylog) Regular() seelog.LoggerInterface {
	return l.Log
}

//得到日志实例
func GetMylog() *mylog {
	return my_log
}
