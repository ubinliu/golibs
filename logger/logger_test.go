package logger

import (
	"fmt"
	"testing"
	"runtime"
)

func TestLogger(t *testing.T){
	Init("appname", "./", LOG_LEVEL_INFO)
	SetLogId("log123456")
	commonFields := map[string]string {
		"ip":"127.0.0.1",
		"uid":"123456",
	}
	WithCommonFields(commonFields)
	Debug("this is debug message %d", 1)
	Info("this is info message %d", 2)
	Notice("this is notice message %d", 3)
	Warn("this is warning message %d", 4)
	Error("this is error message %d", 5)
}

func writeLog(th int, ch chan int){
	i := 100000
	
	for ; i > 0; i-- {
		Info("th %d test Info %d", th, i)
	}
	ch<- th
}

func _TestBenchmakr(t *testing.T){
	cpus := runtime.NumCPU()
	
	fmt.Println("cpu:", cpus)
	runtime.GOMAXPROCS(cpus)

	var i int = 0
	abc := "abcdefghijklmnopqrstuvwxyz"
	message := ""
	for ; i < 50; i++ {
		message += abc
	}
	fmt.Println(message)
	Init("appname", "./", LOG_LEVEL_INFO)
	SetLogId(message)

	ch := make(chan int, cpus)	
	i = 0
	for ; i < cpus; i++ {
		go writeLog(i, ch)
	}

	i = 0
	for ; i < cpus; i++ {
		<-ch
	}
}
