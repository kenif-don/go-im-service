package util

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"encoding/json"
	"sort"
	"strconv"
)

// IndexOfString 查找字符串在数组中的位置
func IndexOfString(target string, arr []string) int {
	sort.Strings(arr)
	index := sort.SearchStrings(arr, target)
	if index < len(arr) && arr[index] == target {
		return index
	}
	return -1
}
func Obj2Obj(m interface{}, obj interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}
func Obj2Str(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func Str2Obj(s string, obj interface{}) error {
	return json.Unmarshal([]byte(s), obj)
}
func Str2Uint64(s string) uint64 {
	i, e := strconv.ParseUint(s, 10, 64)
	if e != nil {
		log.Error(e)
	}
	return i
}
func Uint642Str(i uint64) string {
	return strconv.FormatUint(i, 10)
}
func GetErrMsg(err *utils.Error) string {
	return err.MsgZh
}
