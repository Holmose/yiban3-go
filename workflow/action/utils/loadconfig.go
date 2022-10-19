package actionfunc

import (
	"Yiban3/browser/config"
	"encoding/json"
	"log"
	"os"
)

// LoadSystemConfig 读取系统配置
func LoadSystemConfig(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Panic(err)
	}
	var conf config.ConfigS
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return err
	} else {
		config.CSRF = conf.CSRF
		config.MaxNum = conf.MaxNum
		config.ShowSecond = conf.ShowSecond
		config.CompleteTemplateDelta = conf.CompleteTemplateDelta
		config.MysqlConStr = conf.MysqlConStr
		config.MailUser = conf.MailUser
		config.MailPass = conf.MailPass
		config.MailHost = conf.MailHost
		if conf.SubString != nil {
			config.SubString = conf.SubString
		} else {
			conf.SubString = config.SubString
		}
		config.PerMinute = conf.PerMinute
		config.PerHour = conf.PerHour
		config.WriteSysconf(conf)

		return nil
	}
}
