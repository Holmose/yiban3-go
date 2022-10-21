package schedule

import (
	"Yiban3/Browser/types"
	"log"
	"time"
)

// 将用户浏览器放入浏览器队列中
// 取出一个个人浏览器进行使用
// 取出后将其重新存入浏览器队列中

// 创建一个用户队列
var userQ []browser.User

// 创建一个个人浏览器队列
var browserQ []*browser.Browser

// CreateBrowserByUserQ 根据用户队列创建用户个人浏览器
// CreateBrowserByUserQ 将用户浏览器放入浏览器队列中
func CreateBrowserByUserQ() {
	defer wg.Done()
	count := 0
	for {
		if count > userCount-1 {
			break
		}
		if len(userQ) > 0 {
			user := userQ[0]
			userQ = userQ[1:]

			b := browser.Browser{}
			browser.CreateBrowser(&b, user)
			browserQ = append(browserQ, &b)
			count++
		} else {
			log.Printf("userQ中没有数据，等待中。。。")
			time.Sleep(time.Second * 3)
		}
	}
}
