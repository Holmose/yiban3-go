package clockrpc

import (
	"Yiban3/Flowcharts"
	"errors"
	"log"
)

// Service.Method

type ClockService struct{}

type Args struct {
	Key string
	Iv  string
}

func (ClockService) Clock(args Args, result *string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[err错误：%v]", err)
		}
	}()
	// 合法性判断
	if args.Key != "hFjCM5XBMC6bo3k" && args.Iv != "hONmvJHk" {
		return errors.New("secret key error")
	}
	// 执行打卡
	Flowcharts.ClockOnce()
	*result = "[打卡任务执行结束！]"
	return nil
}

// {"method":"ClockService.Clock","params":[{"Key":"hFjCM5XBMC6bo3k","Iv":"hONmvJHk"}],"id":1}
