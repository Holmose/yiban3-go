package helper

import (
	"log"
	"os"
	"time"
)

var LoggerError *log.Logger
var LoggerInfo *log.Logger

func init() {
	if t, _ := PathExists("log"); !t {
		err := os.Mkdir("log", 0777)
		if err != nil {
			log.Panic("创建文件夹异常")
		}
	}
	s := time.Now().Format("20060102")
	logFile, err := os.OpenFile("log/c_"+s+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常")
	}
	LoggerInfo = log.New(logFile, "Info_", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerError = log.New(logFile, "Error_", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(s string) {
	LoggerInfo.Println(s)
}
func LogError(s string) {
	LoggerError.Println(s)
}
