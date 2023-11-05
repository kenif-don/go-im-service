package util

import (
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
func Map2Obj(m interface{}, obj interface{}) error {
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
	i, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		log.Error(e)
	}
	return uint64(i)
}
