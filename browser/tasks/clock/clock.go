package clock

import (
	"Yiban3/browser/config"
	"Yiban3/browser/tasks/baseaction"
	"Yiban3/browser/types"
	"Yiban3/ecryption/yiban"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GetCompleteByLocation 根据位置返回一条打卡成功的数据
func GetCompleteByLocation(b *browser.Browser, position string) (c browser.CompleteDetail, err error) {
	// 获取3条打卡成功列表
	completeList, err := baseaction.GetCompleteList(b, config.CompleteTemplateDelta)
	if err != nil {
		log.Panicf("用户：%v 获取打卡成功列表失败！", b.User.Username)
		return c, err
	}

	lenCom := len(completeList.Data)
	if lenCom <= 0 {
		return c, errors.New(fmt.Sprintf("用户：%v 没有获取到数据！", b.User.Username))
	}
	completeDetail, err := baseaction.GetMessage(b,
		completeList.Data[rand.Intn(lenCom)])

	if err != nil {
		log.Panicf("用户：%v 获取完成打卡的详细信息失败！", b.User.Username)
		return c, err
	}

	dataJSON := completeDetail.Data.Initiate.FormDataJSON
	var wfName string
	if b.User.IsHoliday {
		wfName = config.SubString["Holiday"]
	} else {
		wfName = config.SubString["Daily"]
	}
	for _, d := range dataJSON {
		if d.Label == "获取定位" {
			address := d.Value.(map[string]interface{})["address"]
			if strings.Contains(completeDetail.Data.WFName, wfName) {
				log.Printf("用户：%v，获取到一条 %v 符合！", b.User.Username, completeDetail.Data.WFName)
				if strings.Contains(address.(string), position) {
					log.Printf("用户：%v，获取到一条 %v 符合！", b.User.Username, address.(string))
					return completeDetail, nil
				} else {
					return c, fmt.Errorf("用户：%v，当前记录 %v 不符合！", b.User.Username, address.(string))
				}
			} else {
				return c, fmt.Errorf("用户：%v，当前记录 %v 不符合！", b.User.Username, address.(string))
			}
		}
	}

	return c, fmt.Errorf("获取指定位置数据失败！")
}

// FillHolidayForm 填写表单 假期
func FillHolidayForm(form browser.FormTask,
	detail browser.Detail,
	completeDetail browser.CompleteDetail) (formEncrypt string, position browser.Position, err error) {
	// 存储表单的ID
	var dataIdM map[string]string
	dataIdM = make(map[string]string)

	for _, component := range form.Data.Form {
		lab := component.Props.Label
		switch {
		case strings.Contains(lab, "今天的体温"):
			dataIdM["TodayTemperatureId"] = component.Id
		case strings.Contains(lab, "身体健康情况"):
			dataIdM["healthStatusID"] = component.Id
		case strings.Contains(lab, "获取定位"):
			dataIdM["positionId"] = component.Id
		case strings.Contains(lab, "家庭所在地区"):
			dataIdM["homePosId"] = component.Id
		case strings.Contains(lab, "目前所在的地区"):
			dataIdM["currentPosId"] = component.Id
		case strings.Contains(lab, "健康码颜色"):
			dataIdM["healCodeColorId"] = component.Id
		case strings.Contains(lab, "正在使用的手机号码"):
			dataIdM["usePhoneId"] = component.Id
		case strings.Contains(lab, "紧急联系人电话"):
			dataIdM["emergencyPhoneId"] = component.Id
		case strings.Contains(lab, "向学校报备的其他情况"):
			dataIdM["otherMessageId"] = component.Id
		}
	}

	nowTime := time.Now().Format("2006-01-02 15:04")

	var formData map[string]interface{}
	formData = make(map[string]interface{})

	// 生成位置信息
	position.Time = nowTime

	for _, dataJSON := range completeDetail.Data.Initiate.FormDataJSON {
		switch {
		case strings.Contains(dataJSON.Label, "获取定位"):
			position.Address = dataJSON.Value.(map[string]interface{})["address"].(string)
			switch dataJSON.Value.(map[string]interface{})["latitude"].(type) {
			case int:
				log.Println("The template FormDataJSON is an int value.")
			case string:
				float, _ := strconv.ParseFloat(dataJSON.Value.(map[string]interface{})["latitude"].(string), 32)
				position.Latitude = float32(float)
				float, _ = strconv.ParseFloat(dataJSON.Value.(map[string]interface{})["longitude"].(string), 32)
				position.Longitude = float32(float)
			case int64:
				log.Println("The template FormDataJSON is an int64 value.")
			default:
				float := dataJSON.Value.(map[string]interface{})["latitude"].(float64)
				position.Latitude = float32(float)
				float, _ = dataJSON.Value.(map[string]interface{})["longitude"].(float64)
				position.Longitude = float32(float)
			}

		case strings.Contains(dataJSON.Label, "家庭所在地区"):
			formData[dataIdM["homePosId"]] = dataJSON.Value.([]interface{})
		case strings.Contains(dataJSON.Label, "目前所在的地区"):
			formData[dataIdM["currentPosId"]] = dataJSON.Value.([]interface{})
		case strings.Contains(dataJSON.Label, "正在使用的手机号码"):
			formData[dataIdM["usePhoneId"]] = dataJSON.Value.(string)
		case strings.Contains(dataJSON.Label, "紧急联系人电话"):
			formData[dataIdM["emergencyPhoneId"]] = dataJSON.Value.(string)
		}
	}

	todayTemperatureList := []string{"36.5", "36.6", "36.7", "36.8"}

	formData[dataIdM["TodayTemperatureId"]] = todayTemperatureList[rand.Intn(len(todayTemperatureList))]
	formData[dataIdM["healthStatusID"]] = "健康"
	formData[dataIdM["positionId"]] = position
	formData[dataIdM["otherMessageId"]] = "无"
	formData[dataIdM["healCodeColorId"]] = "绿"

	formExtend := map[string]interface{}{
		"TaskId": detail.Data.Id,
		"title":  "任务信息",
		"content": []map[string]string{{
			"label": "任务名称",
			"value": detail.Data.Title,
		},
			{
				"label": "发布机构",
				"value": detail.Data.PubOrgName,
			},
		},
	}

	dataJson, err := json.Marshal(formData)
	if err != nil {
		return formEncrypt, position, nil
	}
	extendJson, err := json.Marshal(formExtend)
	if err != nil {
		return formEncrypt, position, nil
	}

	// 拼接
	params := map[string]string{
		"WFId":   detail.Data.WFId,
		"Data":   string(dataJson),
		"Extend": string(extendJson),
	}

	paramsJson, err := json.Marshal(params)
	if err != nil {
		return formEncrypt, position, nil
	}

	formEncrypt, err = yiban.FormEncrypt(string(paramsJson))
	if err != nil {
		return "", position, err
	}

	return formEncrypt, position, nil
}

// FillForm 填写表单 非假期
func FillForm(form browser.FormTask,
	detail browser.Detail,
	completeDetail browser.CompleteDetail) (formEncrypt string, position browser.Position, err error) {
	// 存储表单的ID
	var dataIdM map[string]string
	dataIdM = make(map[string]string)

	var formData map[string]interface{}
	formData = make(map[string]interface{})

	for _, component := range form.Data.Form {
		com := component.Type
		switch {
		case strings.Contains(com, "InputNumber"):
			// 体温
			dataIdM["TodayTemperatureId"] = component.Id
		case strings.Contains(com, "Checkbox"):
			// 健康状态
			dataIdM["healthStatusID"] = component.Id
			formData[dataIdM["healthStatusID"]] = []string{"正常"}
		case strings.Contains(com, "Radio"):
			// 健康状态
			dataIdM["healthStatusID"] = component.Id
			formData[dataIdM["healthStatusID"]] = "正常"
		case strings.Contains(com, "AutoTakePosition"):
			dataIdM["positionId"] = component.Id
		}
	}

	nowTime := time.Now().Format("2006-01-02 15:04")

	position.Time = nowTime

	for _, dataJSON := range completeDetail.Data.Initiate.FormDataJSON {
		if strings.Contains(dataJSON.Label, "获取定位") {
			position.Address = dataJSON.Value.(map[string]interface{})["address"].(string)

			switch dataJSON.Value.(map[string]interface{})["latitude"].(type) {
			case int:
				fmt.Println("is an int value.")
			case string:
				float, _ := strconv.ParseFloat(dataJSON.Value.(map[string]interface{})["latitude"].(string), 32)
				position.Latitude = float32(float)
				float, _ = strconv.ParseFloat(dataJSON.Value.(map[string]interface{})["longitude"].(string), 32)
				position.Longitude = float32(float)
			case int64:
				fmt.Println("is an int64 value.")
			default:
				float := dataJSON.Value.(map[string]interface{})["latitude"].(float64)
				position.Latitude = float32(float)
				float, _ = dataJSON.Value.(map[string]interface{})["longitude"].(float64)
				position.Longitude = float32(float)
			}
		}
	}

	todayTemperatureList := []string{"36.5", "36.6", "36.7", "36.8"}
	lenToday := len(todayTemperatureList)
	if lenToday <= 0 {
		log.Println("体温列表为空!")
	}
	randNum := rand.Intn(lenToday)

	formData[dataIdM["TodayTemperatureId"]] = todayTemperatureList[randNum]
	formData[dataIdM["positionId"]] = position

	formExtend := map[string]interface{}{
		"TaskId": detail.Data.Id,
		"title":  "任务信息",
		"content": []map[string]string{{
			"label": "任务名称",
			"value": detail.Data.Title,
		},
			{
				"label": "发布机构",
				"value": detail.Data.PubOrgName,
			},
		},
	}

	dataJson, err := json.Marshal(formData)
	if err != nil {
		return formEncrypt, position, nil
	}
	extendJson, err := json.Marshal(formExtend)
	if err != nil {
		return formEncrypt, position, nil
	}

	// 拼接
	params := map[string]string{
		"WFId":   detail.Data.WFId,
		"Data":   string(dataJson),
		"Extend": string(extendJson),
	}

	paramsJson, err := json.Marshal(params)
	if err != nil {
		return formEncrypt, position, nil
	}

	formEncrypt, err = yiban.FormEncrypt(string(paramsJson))
	if err != nil {
		return "", position, err
	}

	return formEncrypt, position, nil
}
