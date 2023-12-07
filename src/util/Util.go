package util

import (
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/entity"
	"encoding/binary"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"time"
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
func Str2Float64(s string) float64 {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		log.Error(e)
	}
	return f
}
func Uint642Str(i uint64) string {
	return strconv.FormatUint(i, 10)
}
func float322byte(fs []float32) []byte {
	bytes := make([]byte, 4*len(fs))
	for i, f := range fs {
		bits := math.Float32bits(f)
		binary.LittleEndian.PutUint32(bytes[i*4:(i+1)*4], bits)
	}
	return bytes
}

// GetErrMsg 统一封装的解密失败消息
func GetErrMsg() string {
	msg := &entity.MessageData{
		Type:    1,
		Content: utils.ERR_DECRYPT_FAIL.MsgZh,
	}
	data, e := Obj2Str(msg)
	if e != nil {
		log.Error(utils.ERR_DECRYPT_FAIL)
	}
	return data
}
func GetDecryptingMsg() string {
	msg := &entity.MessageData{
		Type:    1,
		Content: "文件解密中",
	}
	data, e := Obj2Str(msg)
	if e != nil {
		log.Error(utils.ERR_DECRYPT_FAIL)
	}
	return data
}
func Len(str string) int {
	return len([]rune(str))
}
func CurrentTime() uint64 {
	// 获取当前时间戳
	return uint64(time.Now().UnixNano() / 1e6)
}
func CoverMsgData(tp int, content string) (string, error) {
	md := &entity.MessageData{
		Type:    tp,
		Content: content,
	}
	return Obj2Str(md)
}
func CoverSrcData2EnDate(src, dist []byte, beginIndex, endIndex int) []byte {
	diff := endIndex - beginIndex
	res := make([]byte, len(src)-diff+len(dist))
	copy(res[0:beginIndex], src[0:beginIndex])
	copy(res[beginIndex:len(dist)+beginIndex], dist)
	copy(res[len(dist)+beginIndex:], src[endIndex:])
	return res
}
func RevertCoveredData(enData []byte, oldData []byte, beginIndex, endIndex, diff int) []byte {
	res := make([]byte, len(enData)-diff+(endIndex-beginIndex))
	// 将原始数据的前半部分复制到结果中
	copy(res[0:beginIndex], enData[0:beginIndex])
	// 将原始的被替换部分复制到结果中
	copy(res[beginIndex:beginIndex+len(oldData)], oldData)
	// 将原始数据的后半部分复制到结果中
	copy(res[beginIndex+len(oldData):], enData[beginIndex+diff:])
	return res
}
