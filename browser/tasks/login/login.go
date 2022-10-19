package login

import (
	"Yiban3/browser/config"
	"Yiban3/browser/fetcher"
	"Yiban3/browser/types"
	"Yiban3/ecryption/yiban"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Login 使用账号密码进行登录，并返回一个client对象
func Login(b *browser.Browser) ([]byte, error) {
	loginUrl := "https://oauth.yiban.cn/code/html?client_id=95626fa3080300ea&redirect_uri=https://f.yiban.cn/iapp7463"
	postUrl := "https://oauth.yiban.cn/code/usersure"

	user := &b.User
	// 判断是否存在Verify
	if user.Verify != "" {
		// 二次认证
		ret, err := AuthSecond(b, user)
		if err != nil {
			// Verify认证失败，重新登录获取
			log.Printf("[用户：%v, 执行登录指令!]", user.Username)
			goto FirstAuth
		} else {
			log.Printf("[用户：%v, 登录成功！]", user.Username)
		}
		return ret, nil
	}

FirstAuth:
	// 模拟用户登录
	fetch, err := fetcher.Fetch(b.Client, loginUrl)
	if err != nil {
		return nil, err
	} else {
		log.Printf("用户：%v, 获取login页面成功！", user.Username)
	}
	// 获取秘钥
	publicKeyRe := `-----BEGIN PUBLIC KEY-----[\s\S]*-----END PUBLIC KEY-----`
	re := regexp.MustCompile(publicKeyRe)
	matches := re.FindAllString(string(fetch), -1)
	// 加密密码
	encrypt := yiban.LoginEncrypt([]byte(user.Password), []byte(matches[0]))

	// 定义post参数
	params := url.Values{
		"oauth_uname":  []string{user.Username},
		"oauth_upwd":   []string{encrypt},
		"client_id":    []string{"95626fa3080300ea"},
		"redirect_uri": []string{"https://f.yiban.cn/iapp7463"},
		"state":        []string{""},
		"scope":        []string{"1,2,3,4,"},
		"display":      []string{"html"},
	}
	form, err := b.Client.PostForm(postUrl, params)
	if err != nil {
		log.Panic(err)
	} else {
		log.Printf("用户：%v, 提交登录表单成功！", user.Username)
	}
	defer form.Body.Close()

	bodyReader := bufio.NewReader(form.Body)
	ret, err := io.ReadAll(bodyReader)

	if err != nil {
		return nil, err
	}
	if strings.Contains(string(ret), "https:\\/\\/f.yiban.cn\\/iapp7463") {
		log.Printf("用户：%s, 登录成功！\n", user.Username)
	}

	// 获取Verify
	err = GetVerify(b, user)
	if err != nil {
		return nil, err
	}

	return Login(b)
}

func GetVerify(b *browser.Browser, user *browser.User) error {
	// 尝试获取Verify
	iframeUrl := "https://f.yiban.cn/iframe/index?act=iapp7463"
	resp, err := b.Client.Get(iframeUrl)
	if err != nil {
		return err
	}
	if resp.StatusCode == 521 {
		log.Printf("用户：%v, 开始执行Javascript设置Cookies!", user.Username)
	} else {
		log.Printf("用户：%v, 获取Verify成功！", user.Username)
		// 提取Verify
		verifyRe := `verify_request=(.*)&yb_uid`
		re := regexp.MustCompile(verifyRe)
		matches := re.FindAllStringSubmatch(resp.Request.URL.String(), -1)
		if matches == nil {
			return fmt.Errorf("提取Verify失败！")
		}
		user.Verify = matches[0][1]
		return nil
	}
	bodyReader := bufio.NewReader(resp.Body)
	str, err := io.ReadAll(bodyReader)

	if err != nil {
		return err
	}
	// 提取Javascript代码
	jsRe := `.*setTimeout\("([^(]*)\((\d+)\)".*(function.*)</script>`
	re := regexp.MustCompile(jsRe)
	matches := re.FindAllStringSubmatch(string(str), -1)
	callName := matches[0][1]
	callParam := matches[0][2]
	function := matches[0][3]

	replace := strings.Replace(function, `eval("qo=eval;qo(po);");`, `return po;`, -1)

	vm := goja.New()
	_, err = vm.RunString(replace)
	if err != nil {
		return err
	}
	ret := vm.Get(callName)
	callable, ok := goja.AssertFunction(ret)
	if !ok {
		return fmt.Errorf("AssertFunction失败")
	}
	ret, err = callable(goja.Undefined(), vm.ToValue(callParam))
	if err != nil {
		return err
	}
	// 匹配cookie
	cookieRe := `.*cookie='(.*)'`
	re = regexp.MustCompile(cookieRe)
	matches = re.FindAllStringSubmatch(ret.String(), -1)

	// 通过header将字符串转换为cookie
	header := http.Header{}
	header.Add("Cookie", matches[0][1])
	request := http.Request{Header: header}

	cookieURL, err := url.Parse(iframeUrl)
	if err != nil {
		return err
	}
	b.Client.Jar.SetCookies(cookieURL, request.Cookies())
	// 再次执行
	return GetVerify(b, user)
}

func AuthSecond(b *browser.Browser, user *browser.User) ([]byte, error) {
	// 二次认证链接
	authSecondUrl := "https://api.uyiban.com/base/c/auth/yiban?verifyRequest=" +
		user.Verify + "&CSRF=" + config.CSRF

	// 发送请求
	resp, err := b.ClientGet(authSecondUrl)
	//defer resp.Body.Close()

	if err != nil {
		return nil, err
	}
	// 返回结果
	bodyReader := bufio.NewReader(resp.Body)
	bytes, err := io.ReadAll(bodyReader)
	var resJson interface{}
	err = json.Unmarshal(bytes, &resJson)
	if err != nil {
		log.Panic(err)
	}

	strLen := len(string(bytes))

	if strLen > 1000 {
		log.Printf("[用户：%v, Verify认证成功!]", user.Username)
		return bytes, nil
	} else {
		log.Printf("[用户：%v, Verify认证失败!]", user.Username)
		return nil, fmt.Errorf(" Authentication failed")
	}
}
