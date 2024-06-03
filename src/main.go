package main

import (
	"go-im-service/src/configs/conf"
	"go-im-service/src/configs/log"
	"go-im-service/src/util"
	"os"
)

func main() {
	config := &conf.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "console",
		DeviceType: "android",
		ApiHost:    "http://world-master.online:8886",
		WsHost:     "ws://world-master.online:8003",
	}
	conf.InitConfig(config)
	data, e := os.ReadFile("C:\\Users\\ZhuanZ\\Desktop\\QQ截图20240523153100.png")
	if e != nil {
		log.Debug(e)
		return
	}
	url, err := util.Upload2Bunny("123.png", data)
	if err != nil {
		log.Debug(err)
		return
	}
	log.Debug(url)
}
