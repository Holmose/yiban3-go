package mychan

import (
	"sync"
)

type YibanChan struct {
	C    chan interface{} // 用户数据信息通道
	once sync.Once        // 确保只会关闭一次
}

func NewYibanChan() *YibanChan {
	return &YibanChan{
		C: make(chan interface{}, 10),
	}
}
func (uc *YibanChan) SafeClose() {
	uc.once.Do(func() {
		close(uc.C)
	})
}
