package fetcher

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Unicode解码
func ZhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

const (
	userAgent = "yiban"
	origin    = "https://c.uyiban.com"
)

var rateLimiter = time.Tick(time.Second)

func Fetch(client *http.Client, url string) ([]byte, error) {
	<-rateLimiter
	// 添加 Useragent
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", origin)
	resp, err := client.Do(req)
	//defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusAccepted {
		return nil, fmt.Errorf("获取数据失败，存在反爬！%d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)

	all, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	// 解码Unicode
	unicode, err := ZhToUnicode(all)
	if err != nil {
		return nil, err
	}
	return unicode, nil
}
