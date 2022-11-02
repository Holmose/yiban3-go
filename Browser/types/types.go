package browser

import (
	"Yiban3/Browser/config"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Urbaner 易班打卡相关动作
type Urbaner interface {
	GetCompleteList()       // 获取打卡完成的列表 delta: 前几天的数据
	GetDetail()             // 获取表单信息
	GetMessage()            // 获取完成打卡的详细信息
	GetUnCompleteList()     // 获取未打卡的列表
	FetchUnComplete()       // 从未打卡列表中拿出一条
	CreateForm()            //创建打卡表单（获取在服务器已经生成的
	SubmitForm()            // 提交表单
	GetCompleteByLocation() // 根据位置返回一条打卡成功的数据
	FillHolidayForm()       // 填写表单 假期
	FillForm()              // 填写表单 非假期
	Login()                 // 使用账号密码进行登录，并返回一个client对象
}

// User 用户信息结构
type User struct {
	Username string
	Password string
	Verify   string
	// 当前位置
	Position string
	// 邮件接收者
	Mail string
	// 定时任务
	Crontab   string
	Cron      *cron.Cron
	IsHoliday bool
	// 创建时间
	CreateTime string
	// 更新时间
	UpdateTime string
	// 剩余天数
	Day int
}

// Browser 浏览器结构
type Browser struct {
	Client   *http.Client
	Headers  map[string]string
	User     User
	ChanData ChanData
}

// Tasks 打卡任务列表
type Tasks struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []Data `json:"data"`
}
type Data struct {
	TaskID    string `json:"TaskId"`
	State     int    `json:"State"`
	OrgID     string `json:"OrgId"`
	Title     string `json:"Title"`
	Type      int    `json:"Type"`
	StartTime int64  `json:"StartTime"`
	EndTime   int64  `json:"EndTime"`
}

// Detail 粗略信息
type Detail struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data DetData `json:"data"`
}
type DetData struct {
	Id         string `json:"Id"`
	Title      string `json:"Title"`
	WFId       string `json:"WFId"`
	InitiateId string `json:"InitiateId"`
	PubOrgName string `json:"PubOrgName"`
	PubOrgId   string `json:"PubOrgId"`
}

// CompleteDetail 打卡完成的详细信息
type CompleteDetail struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data ComDetData `json:"data"`
}
type Process struct {
	WFID      string        `json:"WFId"`
	Flow      []interface{} `json:"Flow"`
	CCTrigger string        `json:"CCTrigger"`
}
type PersonInfo struct {
	College        string      `json:"College"`
	Profession     string      `json:"Profession"`
	Class          string      `json:"Class"`
	Grade          int         `json:"Grade"`
	Campus         interface{} `json:"Campus"`
	EducationLevel int         `json:"EducationLevel"`
	StudyYear      string      `json:"StudyYear"`
	Number         string      `json:"Number"`
	Name           string      `json:"Name"`
	Gender         int         `json:"Gender"`
	PersonType     string      `json:"PersonType"`
	Mobile         interface{} `json:"Mobile"`
}
type FormDataJSON struct {
	ID        string      `json:"id"`
	Label     string      `json:"label"`
	Value     interface{} `json:"value"`
	Component string      `json:"component"`
}
type Content struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
type ExtendDataJSON struct {
	TaskID  string    `json:"TaskId"`
	Title   string    `json:"title"`
	Content []Content `json:"content"`
}
type Initiate struct {
	ID             string         `json:"Id"`
	SerialNo       string         `json:"SerialNo"`
	UniversityID   string         `json:"UniversityId"`
	PersonID       string         `json:"PersonId"`
	PersonInfo     PersonInfo     `json:"PersonInfo"`
	WFID           string         `json:"WFId"`
	ProcessID      string         `json:"ProcessId"`
	FormDataJSON   []FormDataJSON `json:"FormDataJson"`
	ExtendDataJSON ExtendDataJSON `json:"ExtendDataJson"`
	WorkNode       int            `json:"WorkNode"`
	State          int            `json:"State"`
	StateTime      string         `json:"StateTime"`
	CreateTime     int            `json:"CreateTime"`
}
type ComDetData struct {
	WFName         string        `json:"WFName"`
	Initiate       Initiate      `json:"Initiate"`
	InitiateExtend []interface{} `json:"InitiateExtend"`
}

// FormTask 打卡提交表单
type FormTask struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data FormData `json:"data"`
}
type FormData struct {
	Id     string      `json:"Id"`
	WFName string      `json:"WFName"`
	Form   []Component `json:"Form"`
}
type Component struct {
	Type  string `json:"component"`
	Props Props  `json:"props"`
	Id    string `json:"id"`
}
type Props struct {
	Label    string   `json:"label"`
	Options  []string `json:"options"`
	Required bool     `json:"required"`
	Decimal  bool     `json:"decimal"`
}

// Position 位置结构
type Position struct {
	Address   string  `json:"address"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Time      string  `json:"time"`
}

type ChanData struct {
	FormChan           chan FormTask
	DetailChan         chan Detail
	CompleteDetailChan chan CompleteDetail
	UnCompleteChan     chan Data
}

func CreateBrowser(b *Browser, user User) {
	// 定义client，并使用cookie自动管理cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		// 初始化浏览器失败直接退出
		log.Panic(err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 18 * time.Second,
	}

	// 创建Header信息
	header := map[string]string{
		"Origin":     "https://c.uyiban.com",
		"User-Agent": "yiban",
	}

	b.Client = client
	b.User = user
	b.Headers = header

	chandata := ChanData{
		FormChan:           make(chan FormTask, 2),
		DetailChan:         make(chan Detail, 2),
		CompleteDetailChan: make(chan CompleteDetail, 2),
		UnCompleteChan:     make(chan Data, 2),
	}
	b.ChanData = chandata
}

func (b *Browser) ClientGet(urlString string) (resp *http.Response, err error) {
	// 设置csrf cookie
	header := http.Header{}
	header.Add("Cookie", "csrf_token="+config.CSRF)
	request := http.Request{Header: header}
	cookieURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	b.Client.Jar.SetCookies(cookieURL, request.Cookies())

	// 添加 Useragent
	req, _ := http.NewRequest("GET", urlString, nil)
	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}
	// 发送请求
	resp, err = b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (b *Browser) ClientPost(urlString string, params string) (resp *http.Response, err error) {
	// 设置csrf cookie
	header := http.Header{}
	header.Add("Cookie", "csrf_token="+config.CSRF)
	request := http.Request{Header: header}
	cookieURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	b.Client.Jar.SetCookies(cookieURL, request.Cookies())

	// 添加 Useragent
	req, _ := http.NewRequest("POST", urlString, strings.NewReader(params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	defer req.Body.Close()

	for k, v := range b.Headers {
		req.Header.Set(k, v)
	}
	// 发送请求
	resp, err = b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
