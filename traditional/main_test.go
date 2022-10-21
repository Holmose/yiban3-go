package main_test

import (
	"Yiban3/Browser/config"
	"Yiban3/schedule"
	"log"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	for i := 0; i < 2; i++ {
		timeUnix := time.Now().Unix()
		tm := time.Unix(1661788800, 0)
		log.Printf("程序有效期到 %v", tm.Format("2006年01月02日15:04"))

		if timeUnix <= 1661788800 {
			log.Println("当前有效...")
			schedule.ChanListRunMysql()
			time.Sleep(time.Second * config.ShowSecond)
		} else {
			log.Println("程序已过期，请联系管理员。")
		}
	}
}
