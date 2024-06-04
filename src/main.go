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
	data, e := os.ReadFile("C:\\Users\\ZhuanZ\\Desktop\\微信截图_20240603222353.png")
	if e != nil {
		log.Debug(e)
		return
	}
	url, err := util.Upload2Bunny("234.png", data)
	if err != nil {
		log.Debug(err)
		return
	}
	log.Debug(url)
}
