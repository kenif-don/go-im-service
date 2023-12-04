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
	"github.com/h2non/filetype"
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
	//添加请求头和body
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
		return &resultDTO, nil
	}
	d, e := handlerResult(req, &resultDTO)
	if e != nil {
		return nil, utils.NewError(resultDTO.Code, resultDTO.Msg, resultDTO.Msg)
	}
	resultDTO.Data = d
	return &resultDTO, nil
}
func handlerResult(req *http.Request, resultDTO *dto.ResultDTO) (interface{}, error) {
	if resultDTO.Data == nil {
		return nil, nil
	}
	if IndexOfString(req.URL.Path, conf.Conf.ExUris) == -1 {
		data, e := DecryptAes(resultDTO.Data.(string), conf.Conf.Key)
		if e != nil {
			return nil, e
		}
		return data, nil
	}
	return resultDTO.Data, nil
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
	if conf.GetLoginInfo().User == nil {
		return utils.ERR_NOT_LOGIN
	}
	//参数加密 服务器公钥+自己的私钥 协商出来共享秘钥加密参数
	conf.Conf.Key = SharedAESKey(conf.Conf.Pk, conf.GetLoginInfo().User.PrivateKey, conf.Conf.Prime)
	newData, e := EncryptAes(string(data), conf.Conf.Key)
	if e != nil {
		return log.WithError(e)
	}
	//将字符串赋值给请求对象body
	req.Body = io.NopCloser(bytes.NewBuffer([]byte(newData)))
	req.Header.Add("sign", strings.ToUpper(MD5(sign+newData)))
	return nil
}
func UploadData(data []byte, secret string) (string, *utils.Error) {
	sess, e := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.Conf.Aws.Id, conf.Conf.Aws.Secret, ""),
		Endpoint:         aws.String(conf.Conf.Aws.Endpoint),
		Region:           aws.String(conf.Conf.Aws.Region),
		S3ForcePathStyle: aws.Bool(true),
	})
	if e != nil {
		log.Debug(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	uploader := s3.New(sess)
	//获取文件后缀
	endWith, err := GetFileType(data)
	if err != nil {
		log.Debug(err)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	//文件MD5作为文件名称
	filename := MD5Bytes(data) + "." + endWith
	//将data 加密
	subData := data[3:19]
	subEnData, err := EncryptAes2(subData, secret)
	if err != nil {
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	data = CoverSrcData2EnDate(data, subEnData, 3, 19)
	_, e = uploader.PutObjectWithContext(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(conf.Conf.Aws.Bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(data),
		ACL:    aws.String("public-read"),
	})
	if e != nil {
		log.Error(err)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	// 获取预览URL
	return "https://" + conf.Conf.Aws.Endpoint + "/" + conf.Conf.Aws.Bucket + "/" + filename, nil
}
func Upload(path string, secret string) (string, *utils.Error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Debug(err)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	return UploadData(data, secret)
}
func UploadFile(data []byte, path string, secret string) (string, *utils.Error) {
	if data == nil {
		return Upload(path, secret)
	}
	return UploadData(data, secret)
}
func GetFileType(data []byte) (string, *utils.Error) {
	//取data的前261个
	buffer := data[:261]
	kind, _ := filetype.Match(buffer)
	if kind == filetype.Unknown {
		log.Debug("未知文件类型")
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	//以斜杠分割 取最后一个
	tp := kind.MIME.Value
	tps := strings.Split(tp, "/")
	return tps[len(tps)-1], nil
}
func DownloadFile(url, path string) error {
	// 判断文件是否存在
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	// 创建一个 HTTP 客户端
	client := &http.Client{}
	// 创建一个 GET 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// 关闭响应体
	defer resp.Body.Close()
	// 创建一个文件，用于保存下载的文件
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	// 将响应体的内容写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	// 关闭文件
	file.Close()
	return nil
}
func DecryptFile(path, secret string) ([]byte, *utils.Error) {
	//读取文件
	data, e := os.ReadFile(path)
	if e != nil {
		log.Debug(e)
		return nil, log.WithError(utils.ERR_DECRYPT_FAIL)
	}
	//解密
	subEnData := data[3:35]
	oldData, err := DecryptAes2(subEnData, secret)
	if err != nil {
		return nil, log.WithError(err)
	}
	return RevertCoveredData(data, oldData, 3, 19, len(subEnData)), nil
}
