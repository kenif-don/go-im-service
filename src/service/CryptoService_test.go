package service

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/log"
	"testing"
)

func init() {
	config := &conf.BaseConfig{
		BaseDir:    "..",
		LogSwitch:  "CONSOLE_FILE",
		DeviceType: "2",
		ApiHost:    "http://hp9kwse9.beesnat.com",
		WsHost:     "ws://ggeejj9f.beesnat.com:13191",
	}
	_, err := conf.InitConfig(config)
	if err != nil {
		return
	}
}
func TestCreateDHKey(t *testing.T) {
	//keys := CreateDHKey("262074f1e0e19618f0d2af786779d6ad9e814b", "02")
	//log.Debugf("%+v", keys)
	//{PublicKey:"173c01ee05dafb2d599a52dcb1d2cc992af4f1" PrivateKey:8ef425f457e563bce9bf03fb315d678f2b3f0}
	//155Jq5pu3245d4418M19YnRvau7Rc14hVB2301
	key := SharedAESKey("19311a1a18656914b9381c058c309083022301", "8ef425f457e563bce9bf03fb315d678f2b3f0", "262074f1e0e19618f0d2af786779d6ad9e814b")
	log.Debug(key)
}
