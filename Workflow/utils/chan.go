package utils

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type YibanChan struct {
	C    chan interface{} // 用户数据信息通道
	once sync.Once        // 确保只会关闭一次
}
type CronUser struct {
	UserName   string
	Spec       string
	EntryID    cron.EntryID
	UpdateTime string
}

func NewYibanChan() *YibanChan {
	return &YibanChan{
		C: make(chan interface{}, 10),
	}
}
func (y *YibanChan) SafeClose() {
	y.once.Do(func() {
		close(y.C)
	})
}

// SafeClose 通用关闭
func SafeClose(ch interface{}) (justClosed bool) {

	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	close(ch.(chan interface{}))
	return true
}
