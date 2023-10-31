package util

import (
	"IM-Service/src/configs/conf"
	"IM-Service/src/configs/log"
	"IM-Service/src/dto"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Post 发起POST json请求
func Post(url string, body interface{}) (*dto.ResultDTO, error) {
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", conf.Base.ApiHost+url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	//添加请求头
	addHeader(req, data)
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var resultDTO dto.ResultDTO
	err = json.Unmarshal(result, &resultDTO)
	if err != nil {
		return nil, err
	}
	return &resultDTO, nil
}
func addHeader(req *http.Request, data []byte) {
	req.Header.Add("Content-Type", "application/json")
	//添加签名
	timestamp, sign := GetSign()
	req.Header.Add("timestamp", strconv.FormatInt(timestamp, 10))
	//放行
	if IndexOfString(req.URL.Path, conf.Conf.Data.ExcludeUri) != -1 {
		req.Header.Add("sign", sign)
		return
	}
	param := ""
	if data != nil {
		param = string(data)
	}
	log.Debug(param)
	req.Header.Add("sign", MD5(sign+param))
}
