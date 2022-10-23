package browser

import (
	"sync"
)

type ChanData struct {
	FormChan           chan FormTask
	DetailChan         chan Detail
	CompleteDetailChan chan CompleteDetail
	UnCompleteChan     chan Data
}

type BrowChan interface {
	FormTask | Detail | CompleteDetail | Data
}

type YibanChan struct {
	C    chan interface{} // 用户数据信息通道
	once sync.Once        // 确保只会关闭一次
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
