package util

import (
	"bytes"
	"go-im-service/src/configs/conf"
	utils "go-im-service/src/configs/err"
	"go-im-service/src/configs/log"
	"io/ioutil"
	"net/http"
	"time"
)

//	func Upload2S3(filename string, data []byte) (string, *utils.Error) {
//		sess, e := session.NewSession(&aws.Config{
//			Credentials:      credentials.NewStaticCredentials(conf.Conf.Aws.Id, conf.Conf.Aws.Secret, ""),
//			Endpoint:         aws.String(conf.Conf.Aws.Endpoint),
//			Region:           aws.String(conf.Conf.Aws.Region),
//			S3ForcePathStyle: aws.Bool(true),
//		})
//		if e != nil {
//			log.Debug(e)
//			return "", log.WithError(utils.ERR_UPLOAD_FILE)
//		}
//		uploader := s3.New(sess)
//		_, e = uploader.PutObjectWithContext(context.TODO(), &s3.PutObjectInput{
//			Bucket: aws.String(conf.Conf.Aws.Bucket),
//			Key:    aws.String(filename),
//			Body:   bytes.NewReader(data),
//			ACL:    aws.String("public-read"),
//		})
//		if e != nil {
//			log.Error(e)
//			return "", log.WithError(utils.ERR_UPLOAD_FILE)
//		}
//		// 获取预览URL
//		return "https://" + conf.Conf.Aws.Endpoint + "/" + conf.Conf.Aws.Bucket + "/" + filename, nil
//	}
func Upload2Bunny(filename string, data []byte) (string, *utils.Error) {
	client := &http.Client{
		Timeout: 60 * 60 * time.Second,
	}
	// 创建一个新的请求
	req, e := http.NewRequest("PUT", "https://world-master-put.b-cdn.net/"+conf.Conf.Aws.Id+"/"+filename, bytes.NewReader(data))
	if e != nil {
		log.Error(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	// 设置请求头
	req.Header.Set("User-Agent", "Java-BCDN-Client-1.0.4")
	req.Header.Set("AccessKey", conf.Conf.Aws.Secret)
	req.Header.Set("Accept", "*/*")

	// 发送请求
	resp, e := client.Do(req)
	if e != nil {
		log.Error(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	defer resp.Body.Close()

	// 读取响应体
	_, e = ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Error(e)
		return "", log.WithError(utils.ERR_UPLOAD_FILE)
	}
	// 输出响应
	return conf.Conf.Aws.Endpoint + "/" + filename, nil
}

//func Upload2Bunny(filename string, data []byte) (string, *utils.Error) {
//	cfg := &bunnystorage.Config{
//		StorageZone: conf.Conf.Aws.Id,
//		Key:         conf.Conf.Aws.Secret,
//		Endpoint:    bunnystorage.EndpointSingapore,
//		Timeout:     60 * 60 * time.Second,
//		MaxRetries:  3, //重试次数
//	}
//	client, e := bunnystorage.NewClient(cfg)
//	if e != nil {
//		log.Error(e)
//		return "", log.WithError(utils.ERR_UPLOAD_FILE)
//	}
//	_, e = client.Upload(context.TODO(), "", filename, "", bytes.NewReader(data))
//	if e != nil {
//		log.Error(e)
//		return "", log.WithError(utils.ERR_UPLOAD_FILE)
//	}
//	// 获取预览URL
//	return conf.Conf.Aws.Endpoint + "/" + filename, nil
//}

func Upload2Cos(filename string, data []byte) (string, *utils.Error) {
	return "", nil
}
