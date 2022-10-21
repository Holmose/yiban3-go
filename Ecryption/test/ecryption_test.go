package test

import (
	"Yiban3/Ecryption/yiban"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestAesCbc(t *testing.T) {
	// 解密
	decrtptData := "S01pOTZZdkFEN3VWWTJxTjZPL3VBWm1HeFpKM2lIMEZWV21nWldLbHo2cXRyWFd4V0dSYWhaN0IyYzRUeXAvczlweExKclNMUTUyMzFLbi9WcUg2TWtaOGJDZERvdW9HckNuK1E0bmhwaytBVEhjeWw4MlRJYlJIRzVkZTlYeXUwRkkzREx5YzBJUHg3aEI0RlAwaEtJbGlhR1VKZjY3S3NYZVl0Y2lJUzd2aDBCSzUxTlpDRVQ3VHF3MG5teENmcVpOWnREVWdEUWNsa0ZFNEowUUM2djdiMlljbzVsbzBOYmx1dVRUTnVoTndiTlpyMlZsSDNGTHV6VjRYWlVYSlQyZUhDTFIxN281R1orejRwSnJyRjJSRWFlbDNsSitsQm82YzM3ZHMxVGFHaWZCVHpKUDgwLzhRd216YUtZeHFSZHhxUGVubE5LaVg3Uko4a2dmWTFyTzgwRzlQelJSQkNyWXVNYmx0OXRzcWRkdXNtRHNjcXYrSHRqcUlHUElRTitRNVdJSUdvYjUyT294VkNLa2FrTXhsTzNLcnJ6dzlsVlI1bjA1REY0Y2hSNFVFeldJb1pYbW5EeHVVUzVaTFJWaGZwc0tFSm0rN01nb1A3N1dsQ2MzM01tQmxETjJRVkZKT2lSK040L0lYbDBzY0w3M1M3VHowUGh4Tzh1N1I5TFh3VHhwUGpzRmhzTEJhZ082ZVFlMnVnMllQYzlMd0hYYmlrY0hkRExKcFFyRVVHTlRCU1NjQ3NNUTdwNXFwS0VTMlZka1RDSkV4ZVJYYmdqdnRjSWR0aG9SR3J2cGtmNm5uUnRzQ1JCYUlTVVcxd2RReEVDL2xlTCtteVgvQThyeU9JSkt6QmhzbXZTUXlxTFBKamxkQzM1cDk1ODFodkZoZytBOHpMa2JUL0dMV1VLQURSR2tnRFhKcVdCUXFGbDlxUk5CZUQ2bnZvdXFXVGJqSnFoQXBybzRwVk1CUzRMWmJMWHVDMFhMSUpxb2dqUVBDUTNUNWJscmdJbzJ2S2pPQ1grTDhHeng2OFRraGhEaFVKYllGbnBkQ0tKb0tQTGNzTklYdVJOSnplM3J6QlBtVjVZa0ZLRml2VVVHWmlUZ0VjeVhxaHlIMXcxRGVDVlY1ZHgvVFRRVkdqenArL1VCU095VW9YSVpESzVISElXdGNTd2pNUk56YXp4cHI4SW82M1kxYnk3Nk1Pb1RMa3BYZVZlckwwVW1KWU9CZHNQeTFHUlpzMDduSFJ4R1A2TCtzQ3p6U2d0WDlQN1NmdFo0aE9UTkxTT0xkdFFZeHoreFJ5enppOUFRaG14OG5jbWpSSEtmM3prbDdWR2c2dUZKeWdEQ2UxWTZPWUpsWjdSTnFBQXVvenR6NnkzblArbjUrekZacTNJOVBNNTl3MEU3dC9QcUhiUG89"
	decrypt, err := yiban.FormDecrypt(decrtptData)
	if err != nil {
		t.Error(err)
	}

	// 加密
	encrypt, _ := yiban.FormEncrypt(decrypt)
	if decrtptData != encrypt {
		t.Error("加密与解密方法不对应！")
	}

	// 转为结构体
	par := yiban.Params{}
	json.Unmarshal([]byte(decrypt), &par)

	// 转为JSON
	marshal, _ := json.Marshal(par)

	// 加密JSON
	encrypt, _ = yiban.FormEncrypt(string(marshal))

	// 解密JSON
	formDecrypt, _ := yiban.FormDecrypt(encrypt)

	//JSON中文转Unicode
	b := []byte(formDecrypt)

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		if size == 3 {
			// 中文
			//fmt.Printf("%c %v\n", r, size)
			// 执行替换
			textQuoted := strconv.QuoteToASCII(string(r))
			textUnquoted := textQuoted[1 : len(textQuoted)-1]
			formDecrypt = strings.Replace(formDecrypt, string(r), textUnquoted, 1)
		}

		b = b[size:]

	}

	t.Log(formDecrypt)

	replace := strings.Replace(decrypt, " ", " ", -1)
	t.Log(replace)
	t.Log("请手动格式化进行对比！")
}
