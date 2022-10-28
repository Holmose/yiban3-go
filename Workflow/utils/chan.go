package utils

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
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

type CronUser struct {
	UserName   string
	Spec       string
	UpdateTime string
	entryID    cron.EntryID
}

// PersonalCrons 存储个人定时任务
type PersonalCrons struct {
	infos map[string]CronUser // 个人定时信息
	cron  *cron.Cron
	mutex *sync.Mutex // 同步管理
}

func (p *PersonalCrons) New() {
	p.infos = make(map[string]CronUser)
	var mutex sync.Mutex
	p.mutex = &mutex
	p.cron = cron.New(cron.WithChain())
}
func (p *PersonalCrons) Add(cu CronUser, cmd func()) error {
	_, ok := p.Get(cu.UserName)
	if ok {
		return errors.New(
			fmt.Sprintf("user %v found", cu.UserName))
	}
	p.mutex.Lock()
	id, err := p.cron.AddFunc(cu.Spec, cmd)
	if err != nil {
		p.mutex.Unlock()
		return err
	}
	cu.entryID = id
	p.infos[cu.UserName] = cu
	p.mutex.Unlock()
	return nil
}
func (p *PersonalCrons) Get(username string) (CronUser, bool) {
	p.mutex.Lock()
	val, ok := p.infos[username]
	p.mutex.Unlock()
	return val, ok
}
func (p *PersonalCrons) GetAll() map[string]CronUser {
	return p.infos
}
func (p *PersonalCrons) Remove(username string) error {
	if username == "" {
		p.mutex.Lock()
		for username, cu := range p.infos {
			p.cron.Remove(cu.entryID)
			delete(p.infos, username)
		}
		p.mutex.Unlock()
		return nil
	}
	// 删除任务
	cu, ok := p.Get(username)
	if ok {
		p.mutex.Lock()
		p.cron.Remove(cu.entryID)
		delete(p.infos, username)
		p.mutex.Unlock()
		return nil
	}
	return errors.New(fmt.Sprintf("user %v can't found", username))
}
func (p *PersonalCrons) Len() int {
	p.mutex.Lock()
	length := len(p.infos)
	p.mutex.Unlock()
	return length
}
func (p *PersonalCrons) Start() {
	p.mutex.Lock()
	p.cron.Start()
	p.mutex.Unlock()
}
func (p *PersonalCrons) Stop() {
	p.mutex.Lock()
	p.cron.Stop()
	p.mutex.Unlock()
}

// Union 求数组的并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

// Intersect 求数组的交集
func Intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// Difference 求数组的差集
func Difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
