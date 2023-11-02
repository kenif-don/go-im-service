package util

import (
	"IM-Service/src/configs/conf"
	utils "IM-Service/src/configs/err"
	"IM-Service/src/configs/log"
	"IM-Service/src/dto"
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Post 发起POST json请求
func Post(url string, body interface{}) (*dto.ResultDTO, *utils.Error) {
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", conf.Base.ApiHost+url, nil)
	if err != nil {
		return nil, log.WithError(err)
	}
	//添加请求头
	err = addContent(req, data)
	if err != nil {
		return nil, log.WithError(err)
	}
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, log.WithError(err)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, log.WithError(err)
	}
	var resultDTO dto.ResultDTO
	err = json.Unmarshal(result, &resultDTO)
	if err != nil {
		return nil, log.WithError(err)
	}
	if resultDTO.Code != 200 {
		return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
	}
	return &resultDTO, nil
}
func addContent(req *http.Request, data []byte) error {
	req.Header.Add("Content-Type", "application/json")
	if conf.GetLoginInfo().Token != "" {
		req.Header.Add("v-token", conf.LoginInfo.Token)
	}
	//添加签名
	timestamp, sign := GetSign()
	req.Header.Add("timestamp", strconv.FormatInt(timestamp, 10))
	//放行
	if IndexOfString(req.URL.Path, conf.Conf.ExUris) != -1 || len(data) == 0 {
		req.Header.Add("sign", sign)
		req.Body = io.NopCloser(bytes.NewBuffer(data))
		return nil
	}
	//参数加密 服务器公钥+自己的私钥 协商出来共享秘钥加密参数
	key := SharedAESKey(conf.Conf.Pk, conf.GetLoginInfo().User.PrivateKey, conf.Conf.Prime)
	newData, err := EncryptAes(string(data), key)
	if err != nil {
		return err
	}
	//将字符串赋值给请求对象body
	req.Body = io.NopCloser(bytes.NewBuffer([]byte(newData)))
	log.Debugf("param:%s", newData)
	req.Header.Add("sign", strings.ToUpper(MD5(sign+newData)))
	return nil
}

func Upload(filename string) (string, *utils.Error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.Conf.Aws.Id, conf.Conf.Aws.Secret, ""),
		Endpoint:         aws.String(conf.Conf.Aws.Endpoint),
		Region:           aws.String(conf.Conf.Aws.Region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		log.Debug(err)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Debug(err)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	uploader := s3.New(sess)
	//获取文件后缀
	endWith := filename[strings.LastIndex(filename, ".")+1:]
	//文件MD5作为文件名称
	filename = MD5Bytes(data) + "." + endWith
	_, err = uploader.PutObjectWithContext(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(conf.Conf.Aws.Bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		log.Error(err)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	// 获取预览URL
	req, _ := uploader.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(conf.Conf.Aws.Bucket),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(99 * 12 * 30 * 24 * time.Hour) // 设置URL的有效期限，这里设置为24小时
	if err != nil {
		log.Error(err)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	return url, nil
}
